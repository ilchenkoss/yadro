package xkcd

import (
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
}

func Xkcd(args OutputStruct) {

	scraper.Scrape(args.DatabasePath,
		args.EDBPath,
		args.TempDir,
		args.TempFolderPattern,
		args.TempFilePattern,
		args.ScrapeLimit,
		args.RequestRetries,
		args.Parallel)

}
