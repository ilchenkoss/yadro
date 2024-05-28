package app

import (
	"context"
	"log/slog"
	"myapp/internal-web/adapters/httpserver"
	"myapp/internal-web/adapters/httpserver/handlers"
	"myapp/internal-web/config"
	"os"
	"os/signal"
)

func Run(cfg *config.Config) {

	//main context for interrupt
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	staticHandler := handlers.NewStaticHandler()
	formsHandler := handlers.NewFormsHandler()
	authHandler := handlers.NewAuthHandler(formsHandler)

	//Init Router
	routerHandlers := &httpserver.Handlers{
		AuthHandler:   authHandler,
		StaticHandler: staticHandler,
		FormsHandler:  formsHandler,
	}
	router := httpserver.NewRouter(routerHandlers)

	//init HttpServer
	//httpCtx := context.Background()
	httpServer := httpserver.NewEngine(cfg, router)
	go func() {
		//start httpserver
		slog.Info("Server listening on " + httpServer.Server.Addr)
		httpServerErr := httpServer.Run()
		if httpServerErr != nil {
			slog.Error("Error starting httpServer: ", "error", httpServerErr.Error())
			panic(httpServerErr)
		}
	}()

	<-ctx.Done()
}
