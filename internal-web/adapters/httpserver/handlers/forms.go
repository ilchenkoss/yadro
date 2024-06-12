package handlers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"myapp/internal-web/core/domain"
	"myapp/internal-web/core/port"
	"net/http"
	"time"
)

type FormsHandler struct {
	TemplateExecutor TemplateExecutor
	XkcdAPI          port.XkcdAPI
	AuthHandler      AuthHandler
}

func NewFormsHandler(te TemplateExecutor, xkcdAPI port.XkcdAPI, ah AuthHandler) *FormsHandler {
	return &FormsHandler{
		TemplateExecutor: te,
		XkcdAPI:          xkcdAPI,
		AuthHandler:      ah,
	}
}

func validateToken(t string) error {

	//token validate
	token, _, err := new(jwt.Parser).ParseUnverified(t, jwt.MapClaims{})
	if err != nil {
		return errors.New("token not valid")
	}

	//exp time validate
	expTime, sErr := token.Claims.GetExpirationTime()
	if sErr != nil {
		return errors.New("token not valid")
	}

	uExpTime := expTime.Time

	if time.Now().After(uExpTime) {
		return errors.New("token expired")
	}

	return nil
}

func (fh *FormsHandler) HomeForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.HomeTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")
	//cookie exist
	if cErr == nil {
		//cookies valid
		if cValidErr := c.Valid(); cValidErr != nil {
			fh.AuthHandler.Logout(w, r)
			return
		}
		//token valid
		if vErr := validateToken(c.Value); vErr != nil {
			fh.AuthHandler.Logout(w, r)
			return
		}
		pageData.Logged = true
	}

	eErr := fh.TemplateExecutor.Home(w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (fh *FormsHandler) LoginForm(w http.ResponseWriter, r *http.Request) {

	var pageData = domain.LoginTemplateData{
		Logged: false,
	}

	c, cErr := r.Cookie("access_token")
	//cookie exist
	if cErr == nil {
		//cookies valid
		if cValidErr := c.Valid(); cValidErr != nil {
			fh.AuthHandler.Logout(w, r)
			return
		}
		//token valid
		if vErr := validateToken(c.Value); vErr != nil {
			fh.AuthHandler.Logout(w, r)
			return
		}
		pageData.Logged = true
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
	//cookie !exist
	if cErr != nil {
		eErr := fh.TemplateExecutor.Comics(w, pageData)
		if eErr != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	//cookies valid
	if cValidErr := c.Valid(); cValidErr != nil {
		fh.AuthHandler.Logout(w, r)
		return
	}
	//token valid
	if vErr := validateToken(c.Value); vErr != nil {
		fh.AuthHandler.Logout(w, r)
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

	if len(comics) > 0 {
		pageData.Comics = comics
	}

	eErr := fh.TemplateExecutor.Comics(w, pageData)
	if eErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

}
