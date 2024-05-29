package handlers

import (
	"myapp/internal-web/core/domain"
	"net/http"
)

type FormsHandler struct {
	StaticFS         http.FileSystem
	TemplateExecutor TemplateExecutor
}

func NewFormsHandler(te TemplateExecutor) *FormsHandler {
	return &FormsHandler{
		StaticFS:         http.Dir("./internal-web/storage/static"),
		TemplateExecutor: te,
	}
}

func (sh *FormsHandler) HomeForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.HomeTemplateData{
		Logged: false,
		Login:  "guest",
	}

	c, cErr := r.Cookie("access_token")
	if cErr == nil {
		cValid := c.Valid()
		if cValid == nil {
			pageData.Logged = true
			//pageData.Login = "Friend"
		}
	}
	eErr := sh.TemplateExecutor.Home(&w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
func (sh *FormsHandler) LoginForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.LoginTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")

	if cErr == nil {
		cValid := c.Valid()
		if cValid == nil {
			pageData.Logged = true
		}
	}

	eErr := sh.TemplateExecutor.Login(&w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
func (sh *FormsHandler) ComicsForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.ComicsTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")
	if cErr == nil {
		cValid := c.Valid()
		if cValid == nil {
			pageData.Logged = true
		}
	}

	requestString := r.URL.Query().Get("s")
	if len(requestString) == 0 {
		eErr := sh.TemplateExecutor.Comics(&w, pageData)
		if eErr != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	pageData.Pictures = []string{
		"picture1.jpg",
		"picture1.jpg",
	}

	eErr := sh.TemplateExecutor.Comics(&w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

}
