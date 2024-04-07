package xkcd

import (
	"io"
	"maps"
	"myapp/pkg/database"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	Status_code_comics_end          = 0
	Status_code_good                = 200
	Status_code_response_error      = 1337
	Status_code_response_read_error = 1338
	Status_code_unknown             = 1339
	Status_code_response_parser     = 1340
	Status_code_no_such_host        = 1341
)

var Condition = true //wait comics end or interrupt

func findLastID(data map[int]database.ParsedData) int {

	var maxID int

	for key := range data {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}

func Scrape(dbPath string, scrapeLimit int) database.ScrapeResult {

	//data from db
	dbData := database.ReadDatabase(dbPath)

	//choose ID, where stopped last scrape
	startID := findLastID(dbData.Data) + 1

	//get new data with old
	scrapedData := MainScrape(dbData, scrapeLimit, startID)

	return scrapedData
}

func MainScrape(dbData database.ScrapeResult, scrapeLimit int, startID int) database.ScrapeResult {

	resultData := dbData.Data
	badIDs := maps.Clone(dbData.BadIDs)

	client := http.Client{Timeout: time.Duration(1) * time.Second} //scrape client

	for Condition && scrapeLimit != 0 {

		var data database.ParsedData
		var response bool

		for ID := range dbData.BadIDs {
			if scrapeLimit == 0 || !Condition {
				break
			}
			data, badIDs, response = secondScrape(client, ID, badIDs)
			delete(dbData.BadIDs, ID)
			if response {
				resultData[ID] = data //append data
			}
			scrapeLimit -= 1
		}

		if Condition { //condition need if interrupt, when checking bad IDs

			data, badIDs, response = secondScrape(client, startID, badIDs)

			if response {
				resultData[startID] = data //append data
			}
		}

		scrapeLimit -= 1
		startID += 1
	}

	dbData.Timestamp = time.Now()
	dbData.BadIDs = badIDs
	dbData.Data = resultData

	return dbData
}

func secondScrape(client http.Client, ID int, badIDs map[int]int) (database.ParsedData, map[int]int, bool) {

	url := "https://xkcd.com/" + strconv.Itoa(ID) + "/info.0.json"
	retries := 3

	dataBytes, statusCode := sendRequest(client, url, retries, ID, Status_code_unknown)

	if statusCode != Status_code_good {
		if statusCode != Status_code_comics_end {
			badIDs[ID] = statusCode
		}
		return database.ParsedData{}, badIDs, false
	}

	data, err := responseParser(dataBytes)
	if err != nil {
		badIDs[ID] = Status_code_response_parser
		return database.ParsedData{}, badIDs, false
	}

	return data, badIDs, true

}

func sendRequest(client http.Client, url string, retries int, ID int, status int) ([]byte, int) { // return data, status code

	if retries <= 0 || !Condition { //exit from recursion
		return nil, status
	}
	println("id: " + strconv.Itoa(ID) + ", retries: " + strconv.Itoa(retries))

	resp, err := client.Get(url)

	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			time.Sleep(1 * time.Second)
			return sendRequest(client, url, retries-1, ID, Status_code_no_such_host)
		}
		return sendRequest(client, url, retries-1, ID, Status_code_response_error)
	}

	defer resp.Body.Close()
	//response ok
	if resp.StatusCode == 200 {

		client.Timeout = 1 // reset timeout

		body, errRead := io.ReadAll(resp.Body)

		if errRead != nil {
			return sendRequest(client, url, retries-1, ID, Status_code_response_read_error)
		}

		return body, Status_code_good

	}

	if resp.StatusCode != 404 { //response not ok

		client.Timeout = 5 //add time to response
		return sendRequest(client, url, retries-1, ID, resp.StatusCode)

	}

	if ID == 404 { //funny comics id
		return nil, 200
	}

	if ID != 404 { // if statusCode == 404 and id != 404
		Condition = false //comics end
		return nil, Status_code_comics_end
	}

	return sendRequest(client, url, retries-1, ID, Status_code_unknown)
}
