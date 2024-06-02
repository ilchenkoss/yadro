package config

import (
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	XkcdApiURL   string `yaml:"xkcd_api_url"`
	StaticPath   string `yaml:"static_path"`
	TemplatePath string `yaml:"template_path"`
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

	if decodeErr := yaml.NewDecoder(file).Decode(&config); decodeErr != nil {
		slog.Error("Error decode config file: ", "error", decodeErr.Error())
		return nil, decodeErr
	}

	return &config, nil
}
