package handlers

import (
	"html/template"
	"net/http"
)

type FormsHandler struct {
	StaticFS http.FileSystem
}

func NewFormsHandler() *FormsHandler {
	return &FormsHandler{
		StaticFS: http.Dir("./internal-web/storage/static"),
	}
}

func (sh *FormsHandler) HomeForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./internal-web/storage/template/index.html",
		"./internal-web/storage/template/body/body.html",
		"./internal-web/storage/template/body/main/home.html",
		"./internal-web/storage/template/body/nav/nav.html",
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	exErr := ts.Execute(w, nil)
	if exErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
func (sh *FormsHandler) LoginForm(w http.ResponseWriter, r *http.Request) {

	files := []string{
		"./internal-web/storage/template/index.html",
		"./internal-web/storage/template/body/body.html",
		"./internal-web/storage/template/body/main/login.html",
		"./internal-web/storage/template/body/nav/nav.html",
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	exErr := ts.Execute(w, nil)
	if exErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
func (sh *FormsHandler) ComicsForm(w http.ResponseWriter, r *http.Request) {

	files := []string{
		"./internal-web/storage/template/index.html",
		"./internal-web/storage/template/body/body.html",
		"./internal-web/storage/template/body/main/comics.html",
		"./internal-web/storage/template/body/nav/nav.html",
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	requestString := r.URL.Query().Get("s")
	if len(requestString) == 0 {
		exErr := ts.Execute(w, nil)
		if exErr != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	type PageData struct {
		Pictures []string
	}

	data := PageData{
		Pictures: []string{
			"picture1.jpg",
			"picture1.jpg",
		},
	}

	exErr := ts.Execute(w, data)
	if exErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

}
