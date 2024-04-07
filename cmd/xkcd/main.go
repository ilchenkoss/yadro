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
		DBPath string `yaml:"db_path"`
	} `yaml:"database"`
}

func loadConfig(configPath string) (Config, error) {

	//open file
	file, err := os.Open(configPath)
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

		//add emergency exit
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}()
}

func main() {

	addInterruptHandling()

	//parse flags
	output := flag.Bool("o", false, "output data")
	outputLimit := flag.Int("n", 2, "number of data output")
	configPath := flag.String("c", "config.yaml", "path to config *.yaml file")
	flag.Parse()

	// load config
	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Println("Error config load:", err)
		return
	}

	//check sourceURL
	if config.Scrape.SourceURL == "https://xkcd.com/" {

		args := xkcd.OutputStruct{
			DatabasePath: config.Database.DBPath,
			OutputLimit:  *outputLimit,
			OutputFlag:   *output,
			ScrapeLimit:  config.Scrape.ScrapePagesLimit,
		}

		xkcd.Xkcd(args)

	} else {
		fmt.Printf("Указанный в config.yaml source_url='%s' нельзя обработать.", config.Scrape.SourceURL)
	}

}
