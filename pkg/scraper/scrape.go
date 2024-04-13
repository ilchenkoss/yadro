package scraper

import (
	"fmt"
	"io"
	"myapp/pkg/database"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Condition = true //interrupt
var ParserScore = sync.WaitGroup{}

func findLastID(data map[int]ScrapedData) int {

	var maxID int
	for key := range data {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}

func Scrape(dbPath string,
	eDBPath string,
	tempDir string,
	tempFolderPattern string,
	tempFilePattern string,
	scrapeLimit int,
	requestRetries int,
	parallel int) error {

	//check temp files
	tempFolders := database.FoundTemp(tempDir, tempFolderPattern)
	if len(tempFolders) > 0 {
		//parser and append
		fmt.Println("temp exists")
	}

	//data from db
	dbDataBytes := database.ReadBytesFromFile(dbPath)
	dbData := decodeFileData(dbDataBytes)

	//choose ID, where stopped last scrape
	startID := findLastID(dbData) + 1

	//found missed IDs
	missedIDs := []int{291}

	//scrape new data
	scrapedData := ScrapePuppeteer(parallel, requestRetries, startID, missedIDs, scrapeLimit, dbData)

	//write last data
	scrapedDataBytes := codeFileData(scrapedData)
	dbErr := database.WriteData(dbPath, eDBPath, scrapedDataBytes)

	if dbErr == nil {
		//remove temps
		fmt.Println("all good, db saved")
	} else {
		fmt.Println(dbErr)
	}

	return nil
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
		fmt.Println(scrapeLimit, startID)
		jobs <- startID
		startID++
		scrapeLimit--
	}
	close(finishScrapeCh)
	return
}

func ScrapePuppeteer(parallel int, retries int, startID int, missedIDs []int, scrapeLimit int, dbData map[int]ScrapedData) map[int]ScrapedData {

	// Create buffered channels for jobs and results
	jobs := make(chan int, 1)
	goodScrapesCh := make(chan []byte)
	resultCh := make(chan map[int]ScrapedData, 1)
	finishCh := make(chan struct{})
	finishScrapeCh := make(chan struct{})
	parserEndCh := make(chan struct{})

	// Set scraper WaitGroup
	var swg sync.WaitGroup
	swg.Add(parallel)

	// Create scrape worker goroutines
	for sworker := 1; sworker <= parallel; sworker++ {
		go scrapeWorker(retries, sworker, jobs, goodScrapesCh, &swg, finishScrapeCh)
	}

	// Append temped response

	// Create parser worker
	go parserWorker(dbData, goodScrapesCh, resultCh, finishCh)

	// Append IDs to jobs
	go appendIDs(jobs, scrapeLimit, missedIDs, startID, finishScrapeCh)

	// Launch a goroutine to wait for all jobs to finish
	go func() {
		swg.Wait()
		fmt.Println("here")
		close(jobs)
		go func() {

			ParserScore.Wait()
			close(parserEndCh)

			for {
				select {
				case <-parserEndCh:
					close(goodScrapesCh)
					close(finishCh)
					return
				}
			}
		}()
	}()

	// Process results
	result := make(map[int]ScrapedData)
	for res := range resultCh {
		result = res
		fmt.Println("result")
		close(resultCh)
	}

	return result
}

// worker performs the task on jobs received and sends results to the results channel.
func scrapeWorker(retries int, workerID int, IDsChan chan int, results chan []byte, swg *sync.WaitGroup, finishScrapeCh chan struct{}) {

	//переместить в место создание воркеров
	//create temp folder
	//database.CreateTempFolder(tempDirPath, tempFolderPattern, workerID)
	//remove after wgparser.done() and dbsave.done()

	client := http.Client{Timeout: time.Duration(1) * time.Second} //scrape client
	//we need client for any worker, because we change timeout time

	//for ID := range IDsChan {
	//	fmt.Println("Scrape status: workerID:", workerID, "requestID:", ID)
	//	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", ID)
	//
	//	if !Condition { //interrupt
	//		wg.Done()
	//		return
	//	}
	//
	//	sendRequest(client, url, retries, ID, results)
	//}

	for {
		select {
		case ID := <-IDsChan:
			fmt.Println("Scrape status: workerID:", workerID, "requestID:", ID)
			url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", ID)
			data := sendRequest(client, url, retries, ID)
			if data != nil {
				ParserScore.Add(1)
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

	println("id: " + strconv.Itoa(ID) + ", retries: " + strconv.Itoa(retries))

	resp, err := client.Get(url)

	fmt.Println(resp.StatusCode)

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
		fmt.Println(resp.StatusCode)
		return nil
	}

	if ID == 404 { //funny comics id
		return nil
	}

	if ID != 404 { // if statusCode == 404 and id != 404
		fmt.Println("here??")
		Condition = false //comics end
		return nil
	}

	return nil
}
