package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"myapp/internal-web/core/domain"
	"net/http"
)

type AuthHandler struct {
	ApiURL           string
	TemplateExecutor TemplateExecutor
}

func NewAuthHandler(apiURL string, te TemplateExecutor) *AuthHandler {
	return &AuthHandler{
		ApiURL:           apiURL,
		TemplateExecutor: te,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	pfErr := r.ParseForm()
	if pfErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	type requestLogin struct {
		Login    string
		Password string
	}
	rl := requestLogin{
		Login:    r.Form.Get("auth_login"),
		Password: r.Form.Get("auth_pass"),
	}
	rlJson, jErr := json.Marshal(rl)
	if jErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	rlReader := bytes.NewReader(rlJson)

	type LoginResponse struct {
		Success bool
		Message string
		Token   string
	}

	res, rErr := http.Post(fmt.Sprintf("%s/login", ah.ApiURL), "application/json", rlReader)
	if rErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if res.StatusCode != http.StatusOK {
		data := domain.LoginTemplateData{
			Logged:   false,
			LoginErr: res.Status,
		}
		teErr := ah.TemplateExecutor.Login(&w, data)
		if teErr != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	resBody, rErr := io.ReadAll(res.Body)
	if rErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Println(login, password)
	var resp LoginResponse
	uErr := json.Unmarshal(resBody, &resp)
	if uErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	token := resp.Token

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

	ah.fh.ComicsForm(w, r)
}
