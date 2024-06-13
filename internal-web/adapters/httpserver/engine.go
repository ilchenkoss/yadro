package httpserver

import (
	"context"
	"fmt"
	"myapp/internal-web/config"
	"net/http"
)

type Engine struct {
	Server *http.Server
	Mux    *http.ServeMux
}

func NewEngine(cfg *config.Config, router *http.ServeMux) *Engine {

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

	errServer := engine.Server.ListenAndServe()

	if errServer != nil && errServer != http.ErrServerClosed {
		return fmt.Errorf("Error server: %v", errServer)
	}

	return nil
}

func (engine *Engine) Stop(ctx context.Context) error {
	return engine.Server.Shutdown(ctx)
}
