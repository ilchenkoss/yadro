package xkcd

import (
	"fmt"
	"myapp/pkg/database"
	"sort"
)

type OutputStruct struct {
	DatabasePath string
	OutputLimit  int
	ScrapeLimit  int
}

func prepToPrint(scrapedData database.ScrapeResult, outputLimit int) {
	//print result

	//preparing to print
	keys := make([]int, 0, len(scrapedData.Data))
	toPrint := map[int]database.ParsedData{}

	for k := range scrapedData.Data {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	if len(keys) > outputLimit {
		keys = keys[:outputLimit]
	}

	for _, k := range keys {
		toPrint[k] = scrapedData.Data[k]
	}
	fmt.Println(database.DataToPrint(toPrint))
}

func findLastID(data map[int]database.ParsedData) int {
	if len(data) == 0 {
		return 1
	}
	var maxID int

	for key := range data {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}
func Output(args OutputStruct) {

	scrapedData, err := database.ReadDatabase(args.DatabasePath)
	if err != nil {
		fmt.Println("err")
	}

	//retry for badIDs
	if len(scrapedData.BadIDs) != 0 {
		for ID := range scrapedData.BadIDs {
			scrapedData = MainScrape(scrapedData.Data, scrapedData.BadIDs, 1, ID)
		}
	}

	scrapedData = MainScrape(scrapedData.Data, scrapedData.BadIDs, args.ScrapeLimit, findLastID(scrapedData.Data))
	database.WriteData(args.DatabasePath, scrapedData)

	if args.OutputLimit > 0 {
		prepToPrint(scrapedData, args.OutputLimit)
	}

}
