package xkcd

import (
	"context"
	"encoding/json"
	"myapp/pkg/database"
	"myapp/pkg/indexing"
	"myapp/pkg/scraper"
)

type OutputStruct struct {
	DatabasePath      string
	EDBPath           string
	TempDir           string
	TempFolderPattern string
	TempFilePattern   string

	ScrapeLimit    int
	RequestRetries int
	Parallel       int

	ScrapeCtx       context.Context
	ScrapeCtxCancel context.CancelFunc

	StringRequest string
	IndexPath     string
}

func Xkcd(args OutputStruct) {

	scrapeData, scrapeScore := scraper.Scrape(args.DatabasePath,
		args.EDBPath,
		args.TempDir,
		args.TempFolderPattern,
		args.TempFilePattern,
		args.ScrapeLimit,
		args.RequestRetries,
		args.Parallel,
		args.ScrapeCtx,
		args.ScrapeCtxCancel)

	var indexDB map[string][]int

	if scrapeScore > 0 {
		indexDB = indexing.CreateIndexingDB(scrapeData, args.IndexPath)
	} else {

		indexDBBytes := database.ReadBytesFromFile(args.IndexPath)
		err := json.Unmarshal(indexDBBytes, &indexDB)

		if err != nil || indexDBBytes == nil {
			indexDB = indexing.CreateIndexingDB(scrapeData, args.IndexPath)
		}
	}

	if len(args.StringRequest) > 0 {
		indexing.FindComics(args.StringRequest, indexDB, scrapeData)
	}
}
