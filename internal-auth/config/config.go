package config

import (
	"errors"
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	defaultConfigPath = "./internal-auth/config/config.yaml"
)

type Config struct {
	Env         string        `yaml:"env"`
	StoragePath string        `yaml:"storage_path"`
	TokenTTL    time.Duration `yaml:"token_ttl"`
	Server      ServerConfig  `yaml:"server"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func LoadConfig() (*Config, error) {

	configPath := fetchConfigPath(defaultConfigPath)
	if _, fcpErr := os.Stat(configPath); errors.Is(fcpErr, os.ErrExist) {
		return nil, fcpErr
	}

	var cfg Config

	if rcErr := cleanenv.ReadConfig(configPath, &cfg); rcErr != nil {
		return nil, rcErr
	}

	return &cfg, nil
}

// fetchConfigPath return config file path with priority: flag > env > default
func fetchConfigPath(defaultConfigPath string) string {
	var configPath string

	flag.StringVar(&configPath, "c", "", "path to config file")

	if configPath == "" {
		var exists bool

		configPath, exists = os.LookupEnv("CONFIG_PATH")
		if !exists {
			configPath = defaultConfigPath
		}
		return configPath
	}
	return configPath
}
