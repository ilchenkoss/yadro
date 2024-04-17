package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"myapp/pkg/xkcd"
	"os"
	"os/signal"
)

type Config struct {
	Scrape struct {
		SourceURL        string `yaml:"source_url"`
		ScrapePagesLimit int    `yaml:"scrape_pages_limit"`
		RequestRetries   int    `yaml:"request_retries"`
		Parallel         int    `yaml:"parallel"`
	} `yaml:"scrape"`
	Database struct {
		DBPath            string `yaml:"db_path"`
		TempDir           string `yaml:"temp_dir"`
		TempFolderPattern string `yaml:"temp_folder_pattern"`
		TempFilePattern   string `yaml:"temp_file_pattern"`
	} `yaml:"database"`
}

func loadConfig(configPath string) Config {

	//open file
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error load config:", err)
		//return Config{} //default config??
		panic(err)
	}
	defer file.Close()

	//decode file
	var config Config
	if decodeErr := yaml.NewDecoder(file).Decode(&config); decodeErr != nil {
		fmt.Println("Error load config:", decodeErr)
		//return Config{} //default config??
		panic(decodeErr)
	}

	return config
}

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
	config := loadConfig(*configPath)

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
			StringRequest:   *stringRequest,
		}

		xkcd.Xkcd(args)

	} else {
		fmt.Printf("Указанный в config.yaml source_url='%s' нельзя обработать.", config.Scrape.SourceURL)

	}

}
