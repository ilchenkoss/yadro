package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"myapp/pkg/xkcd"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Scrape struct {
		SourceURL        string `yaml:"source_url"`
		ScrapePagesLimit int    `yaml:"scrape_pages_limit"`
	} `yaml:"scrape"`
	Database struct {
		DBFile string `yaml:"db_file"`
		DBPath string `yaml:"db_path"`
	} `yaml:"database"`
}

func loadConfig() (Config, error) {

	//open file
	file, err := os.Open("config.yaml")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	//decode file
	var config Config
	decoder := yaml.NewDecoder(file)
	if decodeErr := decoder.Decode(&config); decodeErr != nil {
		return Config{}, decodeErr
	}

	return config, nil
}

func addInterruptHandling() {
	sign := make(chan os.Signal, 1)

	//select incoming signals
	signal.Notify(sign, os.Interrupt)

	go func() {
		//wait interrupt
		<-sign
		//change condition
		xkcd.Condition = false
		fmt.Println("Interrupt. Stopping scrape...")
		//add timeout
		timeout := time.After(2 * time.Second)
		go func() {
			select {
			case <-timeout:
				fmt.Println("timeout interrupt")
				os.Exit(1)
			}
		}()
	}()
}

func main() {

	addInterruptHandling()

	// load config
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error config load:", err)
		return
	}

	//parse flags
	output := flag.Bool("o", false, "output data")
	outputLimit := flag.Int("n", 2, "number of data output")
	flag.Parse()

	//check sourceURL
	if config.Scrape.SourceURL == "https://xkcd.com/" {

		if *output {
			outputArgs := xkcd.OutputStruct{
				DatabasePath: config.Database.DBPath + config.Database.DBFile,
				OutputLimit:  *outputLimit,
				ScrapeLimit:  config.Scrape.ScrapePagesLimit,
			}

			xkcd.Output(outputArgs)
		}
	} else {
		fmt.Printf("Указанный в config.yaml source_url='%s' нельзя обработать.", config.Scrape.SourceURL)
	}

}
