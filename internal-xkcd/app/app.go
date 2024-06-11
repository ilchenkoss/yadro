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

	//Repository dependency injection
	comicsRepo := repository.NewComicsRepository(dbConnection)
	userRepo := repository.NewUserRepository(dbConnection)
	weightsRepo := repository.NewWeightsRepository(dbConnection)

	//Service dependency injection
	weightService := service.NewWeightService()
	scraperClient := scraper.NewScraper(1)
	scrapeService := service.NewScrapeService(ctx, scraperClient, cfg.Scrape)
	tokenService := service.NewTokenService(cfg.HttpServer)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, tokenService)

	//init superAdmin
	superAdmin := domain.User{
		Role:     domain.SuperAdmin,
		Login:    os.Getenv("SUPERUSER_LOGIN"),
		Password: os.Getenv("SUPERUSER_PASSWORD"),
	}
	_, csaErr := userService.RegisterSuperAdmin(&superAdmin)
	if csaErr != nil && !errors.Is(csaErr, domain.ErrUserAlreadyExist) {

		if errors.Is(csaErr, domain.ErrPasswordIncorrect) ||
			errors.Is(csaErr, domain.ErrUserNotSuperAdmin) {
			panic("super admin login or password incorrect")
		}

		panic(csaErr)
	}
	slog.Info("SuperAdmin OK")

	//Handlers dependency injection
	limiter := utils.NewLimiter(&cfg.HttpServer)
	fs := util.OSFileSystem{}
	scrapeHandler := handlers.NewScrapeHandler(scrapeService, weightService, comicsRepo, weightsRepo, ctx, cfg, fs)
	searchHandler := handlers.NewSearchHandler(weightsRepo, weightService, comicsRepo, *limiter)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, userRepo)

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
		TokenService:  tokenService,
		UserRepo:      userRepo,
		Limiter:       limiter,
		UserHandler:   userHandler,
		ScrapeHandler: scrapeHandler,
		SearchHandler: searchHandler,
		AuthHandler:   authHandler,
	}
	router := httpserver.NewRouter(routerHandlers)

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