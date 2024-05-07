package app

import (
	"context"
	"fmt"
	"log/slog"
	"myapp/internal/adapters/database"
	"myapp/internal/adapters/httprouter"
	"myapp/internal/adapters/httprouter/handlers"
	"myapp/internal/adapters/scraper"
	"myapp/internal/config"
	"myapp/internal/core/domain"
	"myapp/internal/core/service/httpserver"
	"myapp/internal/core/service/scrape"
	"myapp/internal/core/service/weight"
	"myapp/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run(cfg *config.Config) {

	//main context for interrupt
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	//make migrations
	migrationErr := storage.MakeMigrations(&cfg.Database)
	if migrationErr != nil {
		slog.Error("Error make migrations: ", migrationErr)
		panic(migrationErr)
	}
	slog.Info("Migrations ok")

	//start db connection
	dbConnection, dbConnErr := database.NewConnection(&cfg.Database)
	if dbConnErr != nil {
		panic(dbConnErr)
	}
	pingErr := dbConnection.Ping()
	if pingErr != nil {
		slog.Error("Error ping DB: ", pingErr)
		panic(pingErr)
	}
	slog.Info("Connection to DB ok")

	//insert words positions for weights
	ipErr := dbConnection.InsertPositions(&[]domain.Positions{
		{ID: 0, Position: "transcript"}, {ID: 1, Position: "alt"}, {ID: 2, Position: "title"},
	})
	if ipErr != nil {
		slog.Error("Error insert positions: ", ipErr)
		panic(ipErr)
	}

	//Dependency injection
	weightService := weight.NewWeightService()
	scraperClient := scraper.NewScraper(1)
	scrapeService := scrape.NewScrapeService(ctx, scraperClient, cfg.Scrape)
	scrapeHandler := handlers.NewScrapeHandler(scrapeService, weightService, dbConnection, ctx, cfg)
	searchHandler := handlers.NewSearchHandler(dbConnection, weightService)

	//Init Router
	routerHandlers := &httprouter.Handlers{
		ScrapeHandler: scrapeHandler,
		SearchHandler: searchHandler,
	}
	router := httprouter.NewRouter(routerHandlers)

	//init HttpServer
	httpCtx := context.Background()
	httpServer := httpserver.NewEngine(&cfg.HttpServer, router)
	go func() {
		//start httpserver
		httpServerErr := httpServer.Run()
		if httpServerErr != nil {
			slog.Error("Error starting httpServer: ", httpServerErr)
			panic(httpServerErr)
		}
	}()

	go func() {
		//add auto update every 24h in 3:00AM
		now := time.Now()
		nextUpdate := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
		if now.After(nextUpdate) {
			nextUpdate = nextUpdate.Add(24 * time.Hour)
		}
		timeToNextUpdate := nextUpdate.Sub(now)
		time.Sleep(timeToNextUpdate)

		//add ticker
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				updateClient := http.Client{Timeout: 10 * time.Second}
				resp, _ := updateClient.Get(fmt.Sprintf("http://%s:%s/update", cfg.HttpServer.Host, cfg.HttpServer.Port))
				if resp.StatusCode == http.StatusOK {
					slog.Info("Update successful")
					continue
				}
				slog.Error("Error auto update")
			}
		}
	}()

	<-ctx.Done()

	if ssErr := httpServer.Stop(httpCtx); ssErr != nil {
		slog.Error("Error shutdown http server: ", ssErr)
	}

	if cdbErr := dbConnection.CloseConnection(); cdbErr != nil {
		slog.Error("Error shutdown database", cdbErr)
	}
}
