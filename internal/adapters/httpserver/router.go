package httpserver

import (
	"myapp/internal/adapters/httpserver/handlers"
	"myapp/internal/core/port"
	"net/http"
)

type Handlers struct {
	TokenService  port.TokenService
	UserRepo      port.UserRepository
	ScrapeHandler *handlers.ScrapeHandler
	SearchHandler *handlers.SearchHandler
	AuthHandler   *handlers.AuthHandler
}

func NewRouter(router *Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update", handlers.AdminMiddleware(router.ScrapeHandler.Update, router.TokenService, router.UserRepo))
	mux.HandleFunc("GET /pics", handlers.OrdinaryMiddleware(router.SearchHandler.Search, router.TokenService, router.UserRepo))
	mux.HandleFunc("POST /login", router.AuthHandler.Login)

	return mux
}
