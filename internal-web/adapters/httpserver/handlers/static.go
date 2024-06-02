package handlers

import (
	"net/http"
)

type StaticHandler struct {
	fs http.FileSystem
}

func NewStaticHandler(staticDir string) *StaticHandler {
	return &StaticHandler{
		fs: http.Dir(staticDir)}
}

func (sh *StaticHandler) Static(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/static/", http.FileServer(sh.fs)).ServeHTTP(w, r)
}
