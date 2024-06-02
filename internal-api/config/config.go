package config

import (
	"github.com/joho/godotenv"
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
	DatabasePath string `yaml:"database_path"`
	DatabaseDSN  string `yaml:"database_dsn"`
}

type TempConfig struct {
	TempDir           string `yaml:"temp_dir"`
	TempFolderPattern string `yaml:"temp_folder_pattern"`
	TempFilePattern   string `yaml:"temp_file_pattern"`
}

type HttpServerConfig struct {
	EnvPath          string `yaml:"env_path"`
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	TokenMaxTime     int    `yaml:"token_max_time"`
	ConcurrencyLimit int    `yaml:"concurrency_limit"`
	RateLimit        int    `yaml:"rate_limit"`
}

func GetConfig(configPath string) (*Config, error) {

	var config Config

	file, err := os.Open(configPath)
	defer func(file *os.File) {
		fcErr := file.Close()
		if fcErr != nil {
			//nothing
			return
		}
	}(file)

	if err != nil {
		slog.Error("Error read config file: ", "error", err.Error())
		return nil, err
	}

	//envLoadErr := godotenv.Load("./internal-api/config/.env")
	envLoadErr := godotenv.Load("./internal-api/config/.env.test")
	if envLoadErr != nil {
		slog.Error("Error loading .env file: %v", err)
	}

	if decodeErr := yaml.NewDecoder(file).Decode(&config); decodeErr != nil {
		slog.Error("Error decode config file: ", "error", decodeErr.Error())
		return nil, decodeErr
	}

	envLoadErr := godotenv.Load(config.HttpServer.EnvPath)
	if envLoadErr != nil {
		slog.Error("Error loading .env file: %v", "error", err.Error())
	}

	return &config, nil
}
