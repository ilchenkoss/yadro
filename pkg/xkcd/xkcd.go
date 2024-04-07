package xkcd

import (
	"fmt"
	"myapp/pkg/database"
	"os"
	"sort"
)

type OutputStruct struct {
	DatabasePath string
	OutputFlag   bool
	OutputLimit  int
	ScrapeLimit  int
}

func PrintLimitedData(scrapedData database.ScrapeResult, outputLimit int) {

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

func Xkcd(args OutputStruct) {

	scrapedData := Scrape(args.DatabasePath, args.ScrapeLimit)

	//write last data
	err := database.WriteData(args.DatabasePath, scrapedData)

	if err != nil {

		errSave2 := database.WriteData("temp_db.json", scrapedData)
		if errSave2 != nil {
			fmt.Println("can't save Data")
		} else {
			pwd, _ := os.Getwd()
			fmt.Printf("data saved to %s%s", pwd, "temp_db.json")
		}

	}

	if args.OutputFlag {
		PrintLimitedData(scrapedData, args.OutputLimit)
	}
}
