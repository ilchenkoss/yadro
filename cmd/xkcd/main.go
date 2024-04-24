package main

import (
	"context"
	"flag"
	"fmt"
	"myapp/cmd"
	"myapp/pkg/xkcd"
	"os"
	"os/signal"

	"gopkg.in/yaml.v3"
)

func main() {

	scrapeCtx, scrapeCtxCancel := context.WithCancel(context.Background())

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		//wait interrupt
		<-ctx.Done()
		//change condition
		scrapeCtxCancel()
		fmt.Println("\nInterrupt. Stopping scrape...")
	}()

	//parse flags
	configPath := flag.String("c", "config.yaml", "path to config *.yaml file")
	emergencyDBPath := flag.String("e", "./pkg/database/edb.json", "emergency Database path")
	stringRequest := flag.String("s", "", "string for your request")

	flag.Parse()

	// load config
	config := cmd.GetConfig(*configPath)

	//check sourceURL
	if config.Scrape.SourceURL == "https://xkcd.com/" {

		args := xkcd.OutputStruct{
			DatabasePath:      config.Database.DBPath,
			EDBPath:           *emergencyDBPath,
			TempDir:           config.Database.TempDir,
			TempFolderPattern: config.Database.TempFolderPattern,
			TempFilePattern:   config.Database.TempFilePattern,

			ScrapeLimit:    config.Scrape.ScrapePagesLimit,
			RequestRetries: config.Scrape.RequestRetries,
			Parallel:       config.Scrape.Parallel,

			ScrapeCtx:       scrapeCtx,
			ScrapeCtxCancel: scrapeCtxCancel,

			StringRequest: *stringRequest,
			IndexPath:     config.Database.IndexPath,
		}

		xkcd.Xkcd(args)

	} else {
		fmt.Printf("Указанный в config.yaml source_url='%s' нельзя обработать.", config.Scrape.SourceURL)

	}

}
