package main

import (
	"flag"
	"fmt"
	"myapp/pkg/scraper"
	"myapp/pkg/xkcd"
	"os"
	"os/signal"
	"time"

	"gopkg.in/yaml.v3"
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
		return Config{}
	}

	return config
}

func addInterruptHandling() {
	sign := make(chan os.Signal, 1)

	//select incoming signals
	signal.Notify(sign, os.Interrupt)

	go func() {
		//wait interrupt
		<-sign
		//change condition
		scraper.Condition = false
		fmt.Println("\nInterrupt. Stopping scrape...")

		//add emergency exit
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}()
}

func main() {

	addInterruptHandling()

	//parse flags
	configPath := flag.String("c", "config.yaml", "path to config *.yaml file")
	emergencyDBPath := flag.String("e", "./pkg/database/edb.json", "emergency Database path")

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
		}

		xkcd.Xkcd(args)

	} else {
		fmt.Printf("Указанный в config.yaml source_url='%s' нельзя обработать.", config.Scrape.SourceURL)

	}

}
