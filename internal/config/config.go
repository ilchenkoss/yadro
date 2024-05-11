package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Scrape     ScrapeConfig
	Database   DatabaseConfig
	Temp       TempConfig
	HttpServer HttpServerConfig
}

type ScrapeConfig struct {
	SourceURL        string `yaml:"source_url"`
	ScrapePagesLimit int    `yaml:"scrape_pages_limit"`
	RequestRetries   int    `yaml:"request_retries"`
	Parallel         int    `yaml:"parallel"`
}

type DatabaseConfig struct {
	DatabasePath  string `yaml:"database_path"`
	DatabaseDSN   string `yaml:"database_dsn"`
	MigrationsDSN string `yaml:"migrations_dsn"`
}

type TempConfig struct {
	TempDir           string `yaml:"temp_dir"`
	TempFolderPattern string `yaml:"temp_folder_pattern"`
	TempFilePattern   string `yaml:"temp_file_pattern"`
}

type HttpServerConfig struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	TokenMaxTime     int    `yaml:"token_max_time"`
	ConcurrencyLimit int    `yaml:"concurrency_limit"`
	RateLimit        int    `yaml:"rate_limit"`
}

func GetConfig(configPath string) (*Config, error) {

	var config Config

	file, err := os.Open(configPath)
	defer file.Close()

	if err != nil {
		slog.Error("Error read config file: ", err)
		return nil, err
	}

	if decodeErr := yaml.NewDecoder(file).Decode(&config); decodeErr != nil {
		slog.Error("Error decode config file: ", decodeErr)
		return nil, decodeErr
	}

	return &config, nil
}
