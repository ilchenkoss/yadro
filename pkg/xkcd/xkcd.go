package xkcd

import (
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/scraper"
	"sort"
)

type OutputStruct struct {
	DatabasePath string
	EDBPath      string
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
	fmt.Println(scraper.DataToPrint(toPrint))
}

func Xkcd(args OutputStruct) {

	scrapedData := scraper.Scrape(args.DatabasePath, args.EDBPath, args.ScrapeLimit)

	if args.OutputFlag {
		PrintLimitedData(scrapedData, args.OutputLimit)
	}
}
