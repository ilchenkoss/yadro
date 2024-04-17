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

	scraper.Scrape(args.DatabasePath,
		args.EDBPath,
		args.TempDir,
		args.TempFolderPattern,
		args.TempFilePattern,
		args.ScrapeLimit,
		args.RequestRetries,
		args.Parallel,
		args.ScrapeCtx,
		args.ScrapeCtxCancel)

	if len(args.StringRequest) > 0 {
		indexing.MainIndexing(args.StringRequest, args.DatabasePath, args.IndexPath)
	}
}
