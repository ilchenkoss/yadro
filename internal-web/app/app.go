package app

import (
	"context"
	"log/slog"
	"myapp/internal-web/adapters"
	"myapp/internal-web/adapters/httpserver"
	"myapp/internal-web/adapters/httpserver/handlers"
	"myapp/internal-web/config"
	"os"
	"os/signal"
)

func Run(cfg *config.Config) {

	//main context for interrupt
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	xkcdApi := adapters.NewXkcdAPI(cfg.XkcdApiURL)
	templateExecutor := handlers.NewTemplateExecutor(cfg.TemplatePath)
	staticHandler := handlers.NewStaticHandler(cfg.StaticPath)
	authHandler := handlers.NewAuthHandler(xkcdApi, templateExecutor)
	formsHandler := handlers.NewFormsHandler(templateExecutor, xkcdApi, *authHandler)

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
		slog.Info("Web server listening on " + httpServer.Server.Addr)
		httpServerErr := httpServer.Run()
		if httpServerErr != nil {
			slog.Error("Error starting httpServer: ", "error", httpServerErr.Error())
			panic(httpServerErr)
		}
	}()

	<-ctx.Done()
}
