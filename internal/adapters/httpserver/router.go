package httpserver

import (
	"myapp/internal/adapters/httpserver/handlers"
	"myapp/internal/adapters/httpserver/handlers/utils"
	"myapp/internal/core/port"
	"net/http"
)

type Handlers struct {
	TokenService  port.TokenService
	UserRepo      port.UserRepository
	Limiter       *utils.Limiter
	UserHandler   *handlers.UserHandler
	ScrapeHandler *handlers.ScrapeHandler
	SearchHandler *handlers.SearchHandler
	AuthHandler   *handlers.AuthHandler
}

func NewRouter(router *Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /toadmin", utils.SuperAdminMiddleware(router.UserHandler.ToAdmin, router.TokenService, router.UserRepo, router.Limiter))
	mux.HandleFunc("POST /register", utils.AdminMiddleware(router.UserHandler.Register, router.TokenService, router.UserRepo, router.Limiter))

	mux.HandleFunc("POST /update", utils.AdminMiddleware(router.ScrapeHandler.Update, router.TokenService, router.UserRepo, router.Limiter))
	mux.HandleFunc("GET /pics", utils.OrdinaryMiddleware(router.SearchHandler.Search, router.TokenService, router.UserRepo, router.Limiter))

	mux.HandleFunc("POST /login", utils.GuestMiddleware(router.AuthHandler.Login, router.Limiter))

	return mux
}
