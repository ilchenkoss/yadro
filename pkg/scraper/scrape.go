package scraper

import (
	"context"
	"fmt"
	"io"
	"myapp/pkg/database"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var Condition = true //interrupt or comics end

func findIDs(data map[int]ScrapedData, tempedIDs []int) (int, []int) {
	var maxID int
	var missingIDs []int
	existingIDs := make(map[int]bool)

	for key := range data {
		existingIDs[key] = true
	}

	for _, tempedID := range tempedIDs {
		if !existingIDs[tempedID] {
			missingIDs = append(missingIDs, tempedID)
		}
	}

	for id := range existingIDs {
		if id > maxID {
			maxID = id
		}
	}

	for i := 1; i <= maxID; i++ {
		if !existingIDs[i] && i != 404 { // funny ID
			missingIDs = append(missingIDs, i)
		}
	}

	if 0 > len(missingIDs) && len(missingIDs) < 10 {
		fmt.Printf("Потерянные ID: %v", missingIDs)
	}

	if len(missingIDs) >= 10 {
		fmt.Printf("Missed IDs: %v ...and..more..\n", missingIDs[:10])
	}

	return maxID + 1, missingIDs
}

func Scrape(dbPath string,
	eDBPath string,
	tempDirPath string,
	tempFolderPattern string,
	tempFilePattern string,
	scrapeLimit int,
	requestRetries int,
	parallel int) {

	//data from db
	dbDataBytes := database.ReadBytesFromFile(dbPath)
	dbData := decodeFileData(dbDataBytes)

	//check temp files
	temp := database.FoundTemp(tempDirPath, tempFolderPattern, tempFilePattern)

	//choose ID, where stopped last scrape; get missedIDs
	startID, missedIDs := findIDs(dbData, temp.TempIDs)

	//scrape new data
	startTime := time.Now()
	scrapedData, actualTempPath := ScrapePuppeteer(parallel, requestRetries, startID, missedIDs, scrapeLimit, dbData, tempDirPath, tempFolderPattern, tempFilePattern, temp.TempPaths)
	endTime := time.Now()
	fmt.Printf("Scrape time: %v\n", endTime.Sub(startTime))

	//write last data
	scrapedDataBytes := codeFileData(scrapedData)
	dbErr := database.WriteData(dbPath, eDBPath, scrapedDataBytes)

	if dbErr == nil {
		fmt.Println("Data successfully saved.")
		//remove temps
		for oldTempPath := range temp.TempPaths {
			os.RemoveAll(oldTempPath)
		}
		os.RemoveAll(actualTempPath)
	}
	return
}

func appendIDs(jobs chan int, scrapeLimit int, missedIDs []int, startID int, ctxCancel context.CancelFunc) {
	defer ctxCancel()
	//add missed IDs
	for _, ID := range missedIDs {
		if !Condition || scrapeLimit == 0 {
			break
		}
		jobs <- ID
		scrapeLimit--
	}

	//generate IDs++
	for Condition && scrapeLimit != 0 {
		jobs <- startID
		startID++
		scrapeLimit--
	}
}

func ScrapePuppeteer(parallel int,
	retries int,
	startID int,
	missedIDs []int,
	scrapeLimit int,
	dbData map[int]ScrapedData,
	tempDirPath string,
	tempFolderPattern string,
	tempFilePattern string,
	existedTempFiles map[string][]string) (map[int]ScrapedData, string) {

	// Create buffered channels for jobs and results
	jobs := make(chan int, 1)
	goodScrapesCh := make(chan []byte, 1)
	resultCh := make(chan map[int]ScrapedData, 1)
	finishScrapeCtx, ctxCancel := context.WithCancel(context.Background())

	// Create temp folder
	actualTempFolder := database.CreateTempFolder(tempDirPath, tempFolderPattern)

	// Set scraper WaitGroup
	var swg sync.WaitGroup
	swg.Add(parallel)

	// Set parser WaitGroup
	var pwg sync.WaitGroup

	// Create scrape worker goroutines
	for sworker := 1; sworker <= parallel; sworker++ {
		go scrapeWorker(retries, sworker, jobs, goodScrapesCh, &swg, &pwg, finishScrapeCtx, actualTempFolder, tempFilePattern)
	}

	// Append temped response
	go func() {
		for tempFolder, tempFiles := range existedTempFiles {
			for _, tempFile := range tempFiles {
				pwg.Add(1)
				filePath := fmt.Sprintf("%s/%s", tempFolder, tempFile)
				tempData := database.ReadBytesFromFile(filePath)
				goodScrapesCh <- tempData
			}
		}
	}()

	// Create parser worker
	go parserWorker(dbData, goodScrapesCh, &pwg, resultCh)

	// Append IDs to jobs
	go appendIDs(jobs, scrapeLimit, missedIDs, startID, ctxCancel)

	// Launch a goroutine to wait for all jobs to finish
	go func() {
		swg.Wait()
		close(jobs)
		pwg.Wait()
		close(goodScrapesCh)
	}()

	// Process results
	result := <-resultCh
	close(resultCh)

	return result, actualTempFolder
}

func scrapeWorker(retries int,
	workerID int,
	IDsChan chan int,
	results chan []byte,
	swg *sync.WaitGroup,
	pwg *sync.WaitGroup,
	ctx context.Context,
	actualTempFolder string,
	tempFilePattern string) {

	client := http.Client{Timeout: time.Duration(1) * time.Second}

	for {
		select {
		case ID := <-IDsChan:
			fmt.Printf("Scrape status:\nworkerID: %d, requestID: %d\n", workerID, ID)
			url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", ID)
			data := sendRequest(&client, url, retries, ID)
			if data != nil {
				pwg.Add(1)
				database.SaveTemp(data, actualTempFolder, tempFilePattern, ID)
				results <- data
			}

		case <-ctx.Done():
			swg.Done()
			return
		}
	}
}

func sendRequest(client *http.Client, url string, retries int, ID int) []byte {

	if retries <= 0 { //exit from recursion
		return nil
	}

	resp, err := client.Get(url)

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			time.Sleep(1 * time.Second)
			sendRequest(client, url, retries-1, ID)
			return nil
		}
		sendRequest(client, url, retries-1, ID)
		return nil
	}

	defer resp.Body.Close()

	//response ok
	if resp.StatusCode == http.StatusOK {

		body, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			sendRequest(client, url, retries-1, ID)
			return nil
		}

		return body
	}

	if resp.StatusCode != http.StatusNotFound { //response not ok
		sendRequest(client, url, retries-1, ID)
		return nil
	}

	if ID == 404 { //funny comics id
		return nil
	}

	if ID != 404 { // if statusCode == 404 and id != 404
		Condition = false //comics end
		return nil
	}

	return nil
}
