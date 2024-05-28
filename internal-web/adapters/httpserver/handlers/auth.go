package handlers

import (
	"fmt"
	"net/http"
)

type AuthHandler struct {
	fh *FormsHandler
}

func NewAuthHandler(fh *FormsHandler) *AuthHandler {
	return &AuthHandler{
		fh: fh,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	pfErr := r.ParseForm()
	if pfErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	login := r.Form.Get("auth_login")
	password := r.Form.Get("auth_pass")

	fmt.Println(login, password)

	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "Hello world!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	ah.fh.ComicsForm(w, r)
}
