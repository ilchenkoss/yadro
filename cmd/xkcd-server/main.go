package main

import (
	"flag"
	"myapp/internal-api/app"
	"myapp/internal-api/config"
)

func main() {

	configPath := flag.String("c", "./internal-api/config/config.yaml", "path to cfg *.yaml file")
	flag.Parse()

	cfg, cfgErr := config.GetConfig(*configPath)
	if cfgErr != nil {
		panic(cfgErr)
	}

	app.Run(cfg)
}
