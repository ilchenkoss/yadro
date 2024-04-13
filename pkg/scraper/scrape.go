package scraper

import (
	"fmt"
	"io"
	"myapp/pkg/database"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Condition = true //interrupt

func findIDs(data map[int]ScrapedData) (int, []int) {

	var maxID int
	var IDs []int
	var missingIDs []int

	for key := range data {
		IDs = append(IDs, key)
	}

	sort.Ints(IDs)

	if len(IDs) > 0 {
		maxID = IDs[len(IDs)-1]
	} else {
		maxID = 0
	}

	for i := 1; i < len(IDs); i++ {
		if IDs[i]-IDs[i-1] > 1 {
			for j := IDs[i-1] + 1; j < IDs[i]; j++ {
				missingIDs = append(missingIDs, j)
			}
		}
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
	parallel int) error {

	//check temp files
	tempFolders := database.FoundTemp(tempDirPath, tempFolderPattern)
	if len(tempFolders) > 0 {
		//parser and append
		fmt.Println("temp exists")
	}

	//data from db
	dbDataBytes := database.ReadBytesFromFile(dbPath)
	dbData := decodeFileData(dbDataBytes)

	//choose ID, where stopped last scrape; get missedIDs
	startID, missedIDs := findIDs(dbData)

	//scrape new data
	scrapedData, newTempFolders := ScrapePuppeteer(parallel, requestRetries, startID, missedIDs, scrapeLimit, dbData, tempDirPath, tempFolderPattern)

	//write last data
	scrapedDataBytes := codeFileData(scrapedData)
	dbErr := database.WriteData(dbPath, eDBPath, scrapedDataBytes)

	if dbErr == nil {
		//remove temps
		for oldTempFoldersName := range tempFolders {
			os.Remove(tempDirPath + oldTempFoldersName)
		}
		for _, newTempFoldersName := range newTempFolders {
			fmt.Println(newTempFoldersName)
			//КОСТЫЛЬ
			if len(newTempFoldersName) > 0 {
				os.Remove(newTempFoldersName)
			}
		}
	}
	fmt.Println(dbErr)

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
		fmt.Println("ID", startID)
		startID++
		scrapeLimit--
	}
	close(finishScrapeCh)
	return
}

func ScrapePuppeteer(parallel int, retries int, startID int, missedIDs []int, scrapeLimit int, dbData map[int]ScrapedData, tempDirPath string, tempFolderPattern string) (map[int]ScrapedData, []string) {

	// Create buffered channels for jobs and results
	jobs := make(chan int, 1)
	goodScrapesCh := make(chan []byte, 1)
	resultCh := make(chan map[int]ScrapedData, 1)
	finishScrapeCh := make(chan struct{})

	tempDirsCh := make(chan string, parallel)

	// Set scraper WaitGroup
	var swg sync.WaitGroup
	swg.Add(parallel)

	// Set parser WaitGroup
	var pwg sync.WaitGroup

	// Create scrape worker goroutines
	for sworker := 1; sworker <= parallel; sworker++ {
		go scrapeWorker(retries, sworker, jobs, goodScrapesCh, &swg, &pwg, finishScrapeCh, tempDirPath, tempFolderPattern, tempDirsCh)
	}

	// Append temped response

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

	//костыль
	var tempDirsSlice []string
	go func() {
		for {
			select {
			case tempDir := <-tempDirsCh:
				tempDirsSlice = append(tempDirsSlice, tempDir)
			}

		}
	}()

	// Process results
	result := make(map[int]ScrapedData)

	for res := range resultCh {
		result = res
		close(resultCh)
		close(tempDirsCh)
	}

	return result, tempDirsSlice
}

// worker performs the task on jobs received and sends results to the results channel.
func scrapeWorker(retries int,
	workerID int,
	IDsChan chan int,
	results chan []byte,
	swg *sync.WaitGroup,
	pwg *sync.WaitGroup,
	finishScrapeCh chan struct{},
	tempDirPath string,
	tempFolderPattern string,
	tempFoldersCh chan string) {

	//переместить в место создание воркеров
	//create temp folder
	fmt.Println("create temp")
	tempFolder := database.CreateTempFolder(tempDirPath, tempFolderPattern, workerID)
	fmt.Println("after temp")
	tempFoldersCh <- tempFolder

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
