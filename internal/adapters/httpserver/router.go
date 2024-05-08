package httpserver

import (
	"myapp/internal/adapters/httpserver/handlers"
	"net/http"
)

type Handlers struct {
	ScrapeHandler *handlers.ScrapeHandler
	SearchHandler *handlers.SearchHandler
}

func NewRouter(router *Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update", router.ScrapeHandler.Update)
	mux.HandleFunc("GET /pics", router.SearchHandler.Search)

	return mux
}
