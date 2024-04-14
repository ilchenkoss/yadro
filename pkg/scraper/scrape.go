package scraper

import (
	"fmt"
	"io"
	"myapp/pkg/database"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Condition = true //interrupt

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

	fmt.Println(missingIDs)

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
	scrapedData, newTempFolders := ScrapePuppeteer(parallel, requestRetries, startID, missedIDs, scrapeLimit, dbData, tempDirPath, tempFolderPattern, tempFilePattern, temp.TempPaths)
	endTime := time.Now()
	fmt.Printf("Scrape time: %v\n", endTime.Sub(startTime))

	//100 потоков 6.74

	//write last data
	scrapedDataBytes := codeFileData(scrapedData)
	dbErr := database.WriteData(dbPath, eDBPath, scrapedDataBytes)

	if dbErr == nil {
		fmt.Println("Данные успешно сохранены.")
		//remove temps
		for oldTempFoldersName := range temp.TempPaths {
			os.RemoveAll(oldTempFoldersName)
		}
		for _, newTempFoldersName := range newTempFolders {
			os.RemoveAll(newTempFoldersName)
		}
	}

	return
}

func appendIDs(jobs chan int, scrapeLimit int, missedIDs []int, startID int, finishScrapeCh chan struct{}) {

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
	close(finishScrapeCh)
	return
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
	existedTempFiles map[string][]string) (map[int]ScrapedData, []string) {

	// Create buffered channels for jobs and results
	jobs := make(chan int, 1)
	goodScrapesCh := make(chan []byte, 1)
	resultCh := make(chan map[int]ScrapedData, 1)
	finishScrapeCh := make(chan struct{})

	//create temp folder
	actualTempFolder := database.CreateTempFolder(tempDirPath, tempFolderPattern)

	// Set scraper WaitGroup
	var swg sync.WaitGroup
	swg.Add(parallel)

	// Set parser WaitGroup
	var pwg sync.WaitGroup

	// Create scrape worker goroutines
	for sworker := 1; sworker <= parallel; sworker++ {
		go scrapeWorker(retries, sworker, jobs, goodScrapesCh, &swg, &pwg, finishScrapeCh, actualTempFolder, tempFilePattern)
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
	go appendIDs(jobs, scrapeLimit, missedIDs, startID, finishScrapeCh)

	// Launch a goroutine to wait for all jobs to finish
	go func() {
		swg.Wait()
		close(jobs)
		pwg.Wait()
		close(goodScrapesCh)
	}()

	// Process results
	result := make(map[int]ScrapedData)

	for res := range resultCh {
		result = res
		close(resultCh)
	}

	return result, []string{actualTempFolder}
}

// worker performs the task on jobs received and sends results to the results channel.
func scrapeWorker(retries int,
	workerID int,
	IDsChan chan int,
	results chan []byte,
	swg *sync.WaitGroup,
	pwg *sync.WaitGroup,
	finishScrapeCh chan struct{},
	actualTempFolder string,
	tempFilePattern string) {

	client := http.Client{Timeout: time.Duration(1) * time.Second} //scrape client
	//we need client for any worker, because we change timeout time

	for {
		select {
		case ID := <-IDsChan:
			fmt.Println("Scrape status: workerID:", workerID, "requestID:", ID)
			url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", ID)
			data := sendRequest(client, url, retries, ID)
			if data != nil {
				pwg.Add(1)
				database.SaveTemp(data, actualTempFolder, tempFilePattern, ID)
				results <- data
			}

		case <-finishScrapeCh:
			swg.Done()
			return
		}
	}
}

func sendRequest(client http.Client, url string, retries int, ID int) []byte {

	if retries <= 0 || !Condition { //exit from recursion
		fmt.Println("exit retries end")
		return nil
	}

	fmt.Println("id: " + strconv.Itoa(ID) + ", retries: " + strconv.Itoa(retries))

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

		client.Timeout = 1 // reset timeout
		body, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			sendRequest(client, url, retries-1, ID)
			return nil
		}

		return body
		//need add temp file
	}

	if resp.StatusCode != http.StatusNotFound { //response not ok
		client.Timeout = 5 //add time to response
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
