package xkcd

import (
	"io"
	"myapp/pkg/database"
	"net/http"
	"strconv"
	"time"
)

var Condition = true //wait comics end or interrupt

func MainScrape(scrapedData map[int]database.ParsedData, badIDs map[int]int, scrapeLimit int, startID int) database.ScrapeResult {

	var result = database.ScrapeResult{}

	client := http.Client{Timeout: time.Duration(1) * time.Second} //scrape client

	for Condition && scrapeLimit != 0 {

		var data []byte
		data, badIDs = secondScrape(client, startID, badIDs, 3)

		parsedData, parserErr := responseParser(data)

		if parserErr != nil && Condition == true {
			badIDs[startID] = 1337 //parser error
		} else if Condition == true {
			scrapedData[startID] = parsedData //append data
		}

		scrapeLimit -= 1
		startID += 1
	}

	result.Timestamp = time.Now()
	result.BadIDs = badIDs
	result.Data = scrapedData

	return result
}

func secondScrape(client http.Client, ID int, badIDs map[int]int, retries int) ([]byte, map[int]int) {

	url := "https://xkcd.com/" + strconv.Itoa(ID) + "/info.0.json"
	println("id: " + strconv.Itoa(ID) + ", retries: " + strconv.Itoa(retries))

	if retries > 0 {
		//request
		resp, err := client.Get(url)

		if err != nil {

			badIDs[ID] = 0 //error code

			return nil, badIDs
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
	return nil, badIDs //if retries end
}
