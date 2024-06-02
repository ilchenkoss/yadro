package handlers

import (
	"errors"
	"myapp/internal-web/core/domain"
	"myapp/internal-web/core/port"
	"net/http"
)

type AuthHandler struct {
	xAPI      port.XkcdAPI
	tExecutor TemplateExecutor
}

func NewAuthHandler(xkcdAPI port.XkcdAPI, tExecutor TemplateExecutor) *AuthHandler {
	return &AuthHandler{
		xAPI:      xkcdAPI,
		tExecutor: tExecutor,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	pfErr := r.ParseForm()
	if pfErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	token, lErr := ah.xAPI.Login(r.Form.Get("auth_login"), r.Form.Get("auth_pass"))
	if lErr != nil {
		switch {
		case errors.Is(lErr, domain.ErrUnauthorized):
			data := domain.LoginTemplateData{
				Logged:   false,
				LoginErr: "login or password incorrect",
			}
			teErr := ah.tExecutor.Login(&w, data)
			if teErr != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	cookie := http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/comics", 301)
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", 301)
}
