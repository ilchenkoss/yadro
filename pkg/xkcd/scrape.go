package xkcd

import (
	"io"
	"maps"
	"myapp/pkg/database"
	"net/http"
	"strconv"
	"time"
)

var Condition = true //wait comics end or interrupt

func Scrape(dbPath string, scrapeLimit int) database.ScrapeResult {

	//data from db
	dbData := database.ReadDatabase(dbPath)
	//choose ID, where stopped last scrape
	startID := findLastID(dbData.Data) + 1
	//get new data with old
	scrapedData := MainScrape(dbData, scrapeLimit, startID)

	return scrapedData
}

func findLastID(data map[int]database.ParsedData) int {
	if len(data) == 0 {
		return 0
	}
	var maxID int

	for key := range data {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}

func MainScrape(dbData database.ScrapeResult, scrapeLimit int, startID int) database.ScrapeResult {

	resultData := dbData.Data
	badIDs := maps.Clone(dbData.BadIDs)

	if startID == 1 {
		resultData = map[int]database.ParsedData{} //result of scrape
		badIDs = map[int]int{}                     //container with bad response and errors from scrape
	}

	client := http.Client{Timeout: time.Duration(1) * time.Second} //scrape client

	for Condition && scrapeLimit != 0 {
		var data []byte
		//if bad requests in db, retry to scrape
		if len(dbData.BadIDs) != 0 {
			for ID := range dbData.BadIDs {
				data, badIDs = secondScrape(client, ID, badIDs, 3)
				delete(dbData.BadIDs, ID)
				scrapeLimit -= 1
			}
		} else { //if not bad requests

			data, badIDs = secondScrape(client, startID, badIDs, 3)

			parsedData2, parserErr := responseParser(data)

			if parserErr != nil && Condition == true {
				badIDs[startID] = 1337 //parser error
			} else if Condition == true {
				resultData[startID] = parsedData2 //append data
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

func secondScrape(client http.Client, ID int, badIDs map[int]int, retries int) ([]byte, map[int]int) {

	url := "https://xkcd.com/" + strconv.Itoa(ID) + "/info.0.json"
	println("id: " + strconv.Itoa(ID) + ", retries: " + strconv.Itoa(retries))

	if retries <= 0 {
		return nil, badIDs
	}

	//request
	resp, err := client.Get(url)

	if err != nil {

		badIDs[ID] = 0 //error code
		return secondScrape(client, ID, badIDs, retries-1)
		//return nil, badIDs
	}

	defer resp.Body.Close()

	//response ok
	if resp.StatusCode == 200 {

		client.Timeout = 1 // reset timeout

		body, errRead := io.ReadAll(resp.Body)
		if errRead != nil {
			badIDs[ID] = 13373 //error code of read response
			return secondScrape(client, ID, badIDs, retries-1)
		}
		delete(badIDs, ID)
		return body, badIDs

	} else if resp.StatusCode != 404 { //response not ok

		badIDs[ID] = resp.StatusCode

		client.Timeout = 5 //add time to response
		return secondScrape(client, ID, badIDs, retries-1)

	} else if ID != 404 { // if statusCode == 404 and id != 404

		Condition = false
		return nil, badIDs

	} else { // if funny 404 id

		return nil, badIDs
	}
}
