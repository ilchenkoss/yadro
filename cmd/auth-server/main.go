package main

import (
	"context"
	"log/slog"
	"myapp/internal-auth/app"
	"myapp/internal-auth/config"
	"os"
	"os/signal"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	//main context for interrupt
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg, lcErr := config.LoadConfig()
	if lcErr != nil {
		panic(lcErr)
	}
	log := initLogger(cfg.Env)
	log.Info("Config loaded")

	application := app.New(log, *cfg)
	go func() {
		if arErr := application.AppRun(); arErr != nil {
			panic(arErr)
		}
	}()

	<-ctx.Done()
	application.GRPCStop()
}

func initLogger(envType string) *slog.Logger {
	var logger *slog.Logger
	switch envType {
	case envDev:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	}
	return logger
}
