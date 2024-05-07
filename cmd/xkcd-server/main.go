package main

import (
	"flag"
	"myapp/internal/app"
	"myapp/internal/config"
)

func main() {

	configPath := flag.String("c", "./internal/config/config.yaml", "path to cfg *.yaml file")
	flag.Parse()

	cfg, cfgErr := config.GetConfig(*configPath)
	if cfgErr != nil {
		panic(cfgErr)
	}

	app.Run(cfg)
}
