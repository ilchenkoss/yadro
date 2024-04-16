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

func findExistIDs(data map[int]ScrapedData, tempedIDs []int) map[int]bool {
	existIDs := make(map[int]bool)
	//add db ids to map
	for dbID := range data {
		existIDs[dbID] = true
	}
	//add temp ids to map
	for tempID := range tempedIDs {
		existIDs[tempID] = true
	}
	return existIDs
}

func Scrape(dbPath string,
	eDBPath string,
	tempDirPath string,
	tempFolderPattern string,
	tempFilePattern string,
	scrapeLimit int,
	requestRetries int,
	parallel int,
	scrapeCtx context.Context,
	ScrapeCtxCancel context.CancelFunc) {

	//data from db
	dbDataBytes := database.ReadBytesFromFile(dbPath)
	dbData := decodeFileData(dbDataBytes)

	//check temp files
	temp := database.FoundTemp(tempDirPath, tempFolderPattern, tempFilePattern)

	//choose ID, where stopped last scrape; get existIDs
	existIDs := findExistIDs(dbData, temp.TempIDs)

	//scrape new data
	startTime := time.Now()
	scrapedData, actualTempPath := ScrapePuppeteer(parallel, requestRetries, existIDs, scrapeLimit, dbData, tempDirPath, tempFolderPattern, tempFilePattern, temp.TempPaths, scrapeCtx, ScrapeCtxCancel)
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

func appendIDs(jobs chan int, scrapeLimit int, existIDs map[int]bool, scrapeCtx context.Context) {

	// Generate IDs starting from startID and check against existIDs
	startID := 1
	for {
		select {
		case <-scrapeCtx.Done():
			return
		default:
			if scrapeLimit == 0 {
				return
			}
			if !existIDs[startID] {
				jobs <- startID
				scrapeLimit--
			}
			startID++
		}
	}
}

func ScrapePuppeteer(parallel int,
	retries int,
	existIDs map[int]bool,
	scrapeLimit int,
	dbData map[int]ScrapedData,
	tempDirPath string,
	tempFolderPattern string,
	tempFilePattern string,
	existedTempFiles map[string][]string,
	scrapeCtx context.Context,
	scrapeCtxCancel context.CancelFunc) (map[int]ScrapedData, string) {

	// Create buffered channels for jobs and results
	jobs := make(chan int, 1)
	goodScrapesCh := make(chan []byte, 1)
	resultCh := make(chan map[int]ScrapedData, 1)

	// Create temp folder
	actualTempFolder := database.CreateTempFolder(tempDirPath, tempFolderPattern)

	// Set scraper WaitGroup
	var swg sync.WaitGroup
	swg.Add(parallel)

	// Set parser WaitGroup
	var pwg sync.WaitGroup

	// Create scrape worker goroutines
	for sworker := 1; sworker <= parallel; sworker++ {
		go scrapeWorker(retries, sworker, jobs, goodScrapesCh, &swg, &pwg, scrapeCtx, actualTempFolder, tempFilePattern, scrapeCtxCancel)
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
	go appendIDs(jobs, scrapeLimit, existIDs, scrapeCtx)

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
	tempFilePattern string,
	scrapeCtxCancel context.CancelFunc) {

	client := http.Client{Timeout: time.Duration(1) * time.Second}

	for {
		select {
		case ID := <-IDsChan:
			fmt.Printf("Scrape status:\nworkerID: %d, requestID: %d\n", workerID, ID)
			url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", ID)
			data := sendRequest(&client, url, retries, ID, scrapeCtxCancel)
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

func sendRequest(client *http.Client, url string, retries int, ID int, scrapeCtxCancel context.CancelFunc) []byte {

	if retries <= 0 { //exit from recursion
		return nil
	}

	resp, err := client.Get(url)

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			time.Sleep(1 * time.Second)
			sendRequest(client, url, retries-1, ID, scrapeCtxCancel)
			return nil
		}
		sendRequest(client, url, retries-1, ID, scrapeCtxCancel)
		return nil
	}

	defer resp.Body.Close()

	//response ok
	if resp.StatusCode == http.StatusOK {

		body, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			sendRequest(client, url, retries-1, ID, scrapeCtxCancel)
			return nil
		}

		return body
	}

	if resp.StatusCode != http.StatusNotFound { //response not ok
		sendRequest(client, url, retries-1, ID, scrapeCtxCancel)
		return nil
	}

	if ID == 404 { //funny comics id
		return nil
	}

	if ID != 404 { // if statusCode == 404 and id != 404
		scrapeCtxCancel()
		return nil
	}

	return nil
}
