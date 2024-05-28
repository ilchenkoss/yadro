package main

import (
	"flag"
	"myapp/internal-api/app"
	"myapp/internal-api/config"
	"strings"
)

func main() {

	configPath := flag.String("c", "./internal-api/config/config.yaml", "path to cfg *.yaml file")
	superAdmin := flag.String("sa", "humorist yqS~1v1vKcuMs~", "login password for super admin")
	flag.Parse()

	cfg, cfgErr := config.GetConfig(*configPath)
	if cfgErr != nil {
		panic(cfgErr)
	}

	SuperAdminLoginPassword := strings.Fields(*superAdmin)

	app.Run(cfg, SuperAdminLoginPassword)
}
