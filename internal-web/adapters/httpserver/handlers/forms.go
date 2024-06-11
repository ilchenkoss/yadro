package handlers

import (
	"errors"
	"fmt"
	"myapp/internal-web/core/domain"
	"myapp/internal-web/core/port"
	"net/http"
)

type FormsHandler struct {
	StaticFS         http.FileSystem
	TemplateExecutor TemplateExecutor
	XkcdAPI          port.XkcdAPI
	AuthHandler      AuthHandler
}

func NewFormsHandler(te TemplateExecutor, xkcdAPI port.XkcdAPI, ah AuthHandler, staticPath string) *FormsHandler {
	return &FormsHandler{
		StaticFS:         http.Dir(staticPath),
		TemplateExecutor: te,
		XkcdAPI:          xkcdAPI,
		AuthHandler:      ah,
	}
}

func (sh *FormsHandler) HomeForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.HomeTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")
	if cErr == nil {
		cValid := c.Valid()
		if cValid == nil {
			pageData.Logged = true
		}
	}
	eErr := sh.TemplateExecutor.Home(w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
func (fh *FormsHandler) LoginForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.LoginTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")

	if cErr == nil {
		cValidErr := c.Valid()
		if cValidErr == nil {
			pageData.Logged = true
		}
	}

	eErr := fh.TemplateExecutor.Login(w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (fh *FormsHandler) ComicsForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.ComicsTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")
	cValidErr := c.Valid()

	if cErr != nil && cValidErr != nil {
		eErr := fh.TemplateExecutor.Comics(w, pageData)
		if eErr != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	pageData.Logged = true

	requestString := r.URL.Query().Get("s")
	if len(requestString) == 0 {
		eErr := fh.TemplateExecutor.Comics(w, pageData)
		if eErr != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	descriptionID := r.URL.Query().Get("d")

	comicIndex := r.URL.Query().Get("ci")

	if len(descriptionID) != 0 {
		ucErr := fh.XkcdAPI.UpdateDescription(descriptionID, c.Value)
		if ucErr != nil {
			http.Redirect(w, r, fmt.Sprintf("/comics?s=%s&ci=%s", requestString, comicIndex), 301)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/comics?s=%s&ci=%s", requestString, comicIndex), 301)
		return
	}

	comics, xErr := fh.XkcdAPI.GetComics(requestString, c.Value)

	if xErr != nil {
		switch {
		case errors.Is(xErr, domain.ErrUnauthorized):
			fh.AuthHandler.Logout(w, r)
			return
		case errors.Is(xErr, domain.ErrAuthFailed):
			fh.AuthHandler.Logout(w, r)
			return
		case errors.Is(xErr, domain.ErrToManyRequests):
			pageData.SearchErr = "Превышено количество запросов"
		default:
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	pageData.Comics = comics
	eErr := fh.TemplateExecutor.Comics(w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

}
