package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"myapp/internal-xkcd/adapters/httpserver/handlers/utils"
	"myapp/internal-xkcd/core/domain"
	"myapp/internal-xkcd/core/port"
	"net/http"
)

type AuthHandler struct {
	ac port.AuthClient
}

func NewAuthHandler(ac port.AuthClient) *AuthHandler {
	return &AuthHandler{
		ac,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var reqBody LoginRequest
	dErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if dErr != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if len(reqBody.Login) == 0 || len(reqBody.Password) == 0 {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	tokenString, lErr := ah.ac.Login(reqBody.Login, reqBody.Password)
	if lErr != nil {
		if errors.Is(lErr, domain.ErrUserNotFound) {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}
		http.Error(w, "auth failed", http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(utils.NewLoginResponse(true, "Success", tokenString))
	if err != nil {
		return
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	type RegisterRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var reqBody RegisterRequest
	dErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if dErr != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if len(reqBody.Login) == 0 || len(reqBody.Password) == 0 {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	_, err := ah.ac.Register(reqBody.Login, reqBody.Password, domain.Ordinary)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExist) {
			http.Error(w, "user already exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "user register failed", http.StatusInternalServerError)
		return
	}

	encErr := json.NewEncoder(w).Encode("register user successful")
	if encErr != nil {
		//nothing
		return
	}
}
