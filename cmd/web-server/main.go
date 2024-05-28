package main

import (
	"flag"
	"myapp/internal-web/app"
	"myapp/internal-web/config"
)

func main() {

	configPath := flag.String("c", "./internal-web/config/config.yaml", "path to cfg *.yaml file")
	flag.Parse()

	cfg, cfgErr := config.GetConfig(*configPath)
	if cfgErr != nil {
		panic(cfgErr)
	}

	app.Run(cfg)
}
