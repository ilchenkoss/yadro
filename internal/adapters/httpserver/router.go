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

	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			router.ScrapeHandler.Update(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})
	mux.HandleFunc("/pics", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			router.SearchHandler.Search(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})

	return mux
}
