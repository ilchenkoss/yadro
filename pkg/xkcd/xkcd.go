package xkcd

import (
	"context"
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

	indexDB := indexing.CreateIndexingDB(scrapeData, args.IndexPath)

	if len(args.StringRequest) > 0 {
		indexing.ReturnComics(args.StringRequest, indexDB, scrapeData)
	}
}
