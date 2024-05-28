package httpserver

import (
	"myapp/internal-web/adapters/httpserver/handlers"
	"net/http"
)

type Handlers struct {
	*handlers.AuthHandler
	*handlers.StaticHandler
	*handlers.FormsHandler
}

func NewRouter(router *Handlers) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /static/*", router.StaticHandler.Static)

	mux.HandleFunc("POST /login", router.AuthHandler.Login)

	mux.HandleFunc("GET /", router.FormsHandler.HomeForm)
	mux.HandleFunc("GET /login", router.FormsHandler.LoginForm)
	mux.HandleFunc("GET /comics", router.FormsHandler.ComicsForm)

	return mux
}
