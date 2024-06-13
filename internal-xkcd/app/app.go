package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"myapp/internal-xkcd/adapters"
	"myapp/internal-xkcd/adapters/database"
	"myapp/internal-xkcd/adapters/database/repository"
	"myapp/internal-xkcd/adapters/grpc/auth"
	user "myapp/internal-xkcd/adapters/grpc/user"
	"myapp/internal-xkcd/adapters/httpserver"
	"myapp/internal-xkcd/adapters/httpserver/handlers"
	"myapp/internal-xkcd/adapters/httpserver/handlers/utils"
	"myapp/internal-xkcd/adapters/scraper"
	"myapp/internal-xkcd/config"
	"myapp/internal-xkcd/core/domain"
	"myapp/internal-xkcd/core/service"
	"myapp/internal-xkcd/core/util"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run(cfg *config.Config) {

	//main context for interrupt
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	//start db connection
	dbConnection, dbConnErr := database.NewConnection(&cfg.Database)
	if dbConnErr != nil {
		panic(dbConnErr)
	}
	pingErr := dbConnection.Ping()
	if pingErr != nil {
		slog.Error("Error ping DB: ", "error", pingErr.Error())
		panic(pingErr)
	}
	slog.Info("Connection to DB ok")

	//make migrations
	migrationErr := dbConnection.MakeMigrations()
	if migrationErr != nil {
		slog.Error("Error make migrations: ", "error", migrationErr.Error())
		panic(migrationErr)
	}
	slog.Info("Migrations ok")

	authClient, aErr := auth.NewAuth(&cfg.AuthGRPC, ctx)
	if aErr != nil {
		panic(aErr)
	}
	userClient, uErr := user.NewUser(&cfg.AuthGRPC, ctx)
	if uErr != nil {
		panic(uErr)
	}

	comicsRepo := repository.NewComicsRepository(dbConnection)
	weightsRepo := repository.NewWeightsRepository(dbConnection)

	//Service dependency injection
	weightService := service.NewWeightService()
	scraperClient := scraper.NewScraper(1)
	scrapeService := service.NewScrapeService(ctx, scraperClient, cfg.Scrape)

	//init superAdmin
	_, csaErr := authClient.Register(os.Getenv("SUPERUSER_LOGIN"), os.Getenv("SUPERUSER_PASSWORD"), domain.SuperUser)
	if csaErr != nil {
		if !errors.Is(csaErr, domain.ErrUserAlreadyExist) {
			panic(csaErr)
		}
	}
	slog.Info("SuperAdmin OK")

	//Handlers dependency injection
	gptAPI := adapters.NewGptAPI()
	limiter := utils.NewLimiter(&cfg.HttpServer)
	fs := util.OSFileSystem{}
	authHandler := handlers.NewAuthHandler(authClient)
	userHandler := handlers.NewUserHandler(userClient)
	scrapeHandler := handlers.NewScrapeHandler(scrapeService, weightService, comicsRepo, weightsRepo, ctx, cfg, fs)
	searchHandler := handlers.NewSearchHandler(weightsRepo, weightService, comicsRepo, gptAPI, *limiter)

	//insert words positions for weights if not exist
	ipErr := weightsRepo.InsertPositions(&[]domain.Positions{
		{ID: 0, Position: "transcript"}, {ID: 1, Position: "alt"}, {ID: 2, Position: "title"},
	})
	if ipErr != nil {
		slog.Error("Error insert positions: ", "error", ipErr.Error())
		panic(ipErr)
	}

	//Init Router
	routerHandlers := &httpserver.Handlers{
		Limiter:       limiter,
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		ScrapeHandler: scrapeHandler,
		SearchHandler: searchHandler,
	}
	router := httpserver.NewRouter(routerHandlers, authClient, userClient)

	//init HttpServer
	httpCtx := context.Background()
	httpServer := httpserver.NewEngine(&cfg.HttpServer, router)
	go func() {
		//start httpserver
		slog.Info("API server listening on " + httpServer.Server.Addr)
		httpServerErr := httpServer.Run()
		if httpServerErr != nil {
			slog.Error("Error starting httpServer: ", "error", httpServerErr.Error())
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
		slog.Error("Error shutdown http server: ", "error", ssErr.Error())
	}

	if cdbErr := dbConnection.CloseConnection(); cdbErr != nil {
		slog.Error("Error shutdown database", "error", cdbErr.Error())
	}
}
