package xkcd

import (
	"myapp/pkg/database"
	"testing"
)

//
//type Config struct {
//	Scrape struct {
//		SourceURL        string `yaml:"source_url"`
//		ScrapePagesLimit int    `yaml:"scrape_pages_limit"`
//	} `yaml:"scrape"`
//	Database struct {
//		DBFile string `yaml:"db_file"`
//		DBPath string `yaml:"db_path"`
//	} `yaml:"database"`
//}
//
//func loadConfig() (Config, error) {
//
//	//open file
//	file, err := os.Open("config.yaml")
//	if err != nil {
//		return Config{}, err
//	}
//	defer file.Close()
//
//	//decode file
//	var config Config
//	decoder := yaml.NewDecoder(file)
//	if decodeErr := decoder.Decode(&config); decodeErr != nil {
//		return Config{}, decodeErr
//	}
//
//	return config, nil
//}

func TestMissID(t *testing.T) {

	scrapeIDs := 10

	scrapeResult := MainScrape(database.ScrapeResult{}, scrapeIDs, 1)

	var IDs []int
	for goodID := range scrapeResult.Data {
		IDs = append(IDs, goodID)
	}
	for badID := range scrapeResult.BadIDs {
		IDs = append(IDs, badID)
	}

	if len(IDs) != scrapeIDs {

		t.Errorf("\nResult was incorrect. \n scrapes: %d, \n IDs: %d.", scrapeIDs, len(IDs))
	}
}
