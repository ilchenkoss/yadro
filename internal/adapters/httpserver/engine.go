package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"myapp/internal/config"
	"net/http"
)

type Engine struct {
	Server *http.Server
	Mux    *http.ServeMux
}

func NewEngine(cfg *config.HttpServerConfig, router *http.ServeMux) *Engine {

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: router,
	}

	return &Engine{
		Server: server,
		Mux:    router,
	}
}

func (engine *Engine) Run() error {

	slog.Info("Server listening on " + engine.Server.Addr)
	errServer := engine.Server.ListenAndServe()

	if errServer != nil && errServer != http.ErrServerClosed {
		return fmt.Errorf("Error server: %v", errServer)
	}

	return nil
}

func (engine *Engine) Stop(ctx context.Context) error {
	return engine.Server.Shutdown(ctx)
}
