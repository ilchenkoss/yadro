package main

import (
	"flag"
	"myapp/internal-xkcd/app"
	"myapp/internal-xkcd/config"
)

func main() {

	configPath := flag.String("c", "./internal-xkcd/config/config.yaml", "path to cfg *.yaml file")
	flag.Parse()

	cfg, cfgErr := config.GetConfig(*configPath)
	if cfgErr != nil {
		panic(cfgErr)
	}

	app.Run(cfg)
}
