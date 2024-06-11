package httpserver

import (
	"myapp/internal-xkcd/adapters/httpserver/handlers"
	"myapp/internal-xkcd/adapters/httpserver/handlers/utils"
	"myapp/internal-xkcd/core/port"
	"net/http"
)

type Handlers struct {
	Limiter       *utils.Limiter
	ScrapeHandler *handlers.ScrapeHandler
	SearchHandler *handlers.SearchHandler
	AuthHandler   *handlers.AuthHandler
}

func NewRouter(router *Handlers, ac port.AuthClient) *http.ServeMux {
	mux := http.NewServeMux()

	//mux.HandleFunc("POST /toadmin", utils.SuperUserMiddleware(router.UserHandler.ToAdmin, ac, router.Limiter))
	mux.HandleFunc("POST /register", utils.AdminMiddleware(router.AuthHandler.Register, ac, router.Limiter))

	mux.HandleFunc("POST /update", utils.AdminMiddleware(router.ScrapeHandler.Update, ac, router.Limiter))
	mux.HandleFunc("GET /pics", utils.OrdinaryMiddleware(router.SearchHandler.Search, ac, router.Limiter))
	mux.HandleFunc("GET /desc", utils.OrdinaryMiddleware(router.SearchHandler.Description, ac, router.Limiter))

	mux.HandleFunc("POST /login", utils.GuestMiddleware(router.AuthHandler.Login, router.Limiter))

	return mux
}
