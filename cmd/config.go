package cmd

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Scrape   ScrapeConfig
	Database DatabaseConfig
}

type ScrapeConfig struct {
	SourceURL        string `yaml:"source_url"`
	ScrapePagesLimit int    `yaml:"scrape_pages_limit"`
	RequestRetries   int    `yaml:"request_retries"`
	Parallel         int    `yaml:"parallel"`
}

type DatabaseConfig struct {
	DBPath            string `yaml:"db_path"`
	IndexPath         string `yaml:"index_path"`
	TempDir           string `yaml:"temp_dir"`
	TempFolderPattern string `yaml:"temp_folder_pattern"`
	TempFilePattern   string `yaml:"temp_file_pattern"`
}

func getDefaultConfig() *Config {
	slog.Info("Used default Config")
	return &Config{
		Scrape: ScrapeConfig{
			SourceURL:        "https://xkcd.com/",
			ScrapePagesLimit: -1,
			RequestRetries:   3,
			Parallel:         50},
		Database: DatabaseConfig{DBPath: "./pkg/database/database.json",
			IndexPath:         "./pkg/database/index.json",
			TempDir:           "./pkg/database/",
			TempFolderPattern: "temp_xkcd_",
			TempFilePattern:   "response_xkcd"},
	}
}

func GetConfig(configPath string) Config {

	var config Config

	file, err := os.Open(configPath)
	defer file.Close()

	if err != nil {
		return *getDefaultConfig()
	}

	if decodeErr := yaml.NewDecoder(file).Decode(&config); decodeErr != nil {
		return *getDefaultConfig()
	}

	return config
}
