package xkcd

import (
	"fmt"
	"myapp/pkg/database"
	"os"
	"sort"
)

type OutputStruct struct {
	DatabasePath string
	OutputLimit  int
	ScrapeLimit  int
}

func Output(args OutputStruct) {

	var startID int
	var scrapedData database.ScrapeResult
	parsedData := map[int]database.ParsedData{} //result of scrape
	badIDs := map[int]int{}                     //container with bad response and errors from scrape

	dataBytes, databaseErr := database.ReadData(args.DatabasePath) //data

	if databaseErr == nil { //db read ok

		dbData, decodeErr := database.DecodeData(dataBytes) //try decode

		if decodeErr != nil { //decode err
			panic(decodeErr)

		} else { //decode good

			parsedData = dbData.Data
			badIDs = dbData.BadIDs

			//start from last + 1
			startID = findLastID(dbData.Data) + 1
		}

	} else { //error reading database

		if os.IsNotExist(databaseErr) { //if database not exists
			startID = 1
		} else {
			panic(databaseErr) //read error and file exists
		}
	}

	//retry for badIDs
	if len(badIDs) != 0 {
		for ID := range badIDs {
			scrapedData = MainScrape(parsedData, badIDs, 1, ID)
		}
	}

	scrapedData = MainScrape(parsedData, badIDs, args.ScrapeLimit, startID)
	database.WriteData(args.DatabasePath, scrapedData)

	if args.OutputLimit > 0 { //print result

		//preparing to print
		keys := make([]int, 0, len(scrapedData.Data))
		toPrint := map[int]database.ParsedData{}

		for k := range scrapedData.Data {
			keys = append(keys, k)
		}

		sort.Ints(keys)

		if len(keys) > args.OutputLimit {
			keys = keys[:args.OutputLimit]
		}

		// Выводим значения в порядке отсортированных ключей
		for _, k := range keys {
			toPrint[k] = scrapedData.Data[k]
		}
		fmt.Println(database.DataToPrint(toPrint))
	}

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
