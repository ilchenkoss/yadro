package handlers

import (
	"encoding/json"
	"errors"
	"myapp/internal-api/adapters/httpserver/handlers/utils"
	"myapp/internal-api/core/domain"
	"myapp/internal-api/core/port"
	"net/http"
)

type AuthHandler struct {
	svc port.AuthService
}

func NewAuthHandler(svc port.AuthService) *AuthHandler {
	return &AuthHandler{
		svc,
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

	tokenString, lErr := ah.svc.Login(reqBody.Login, reqBody.Password)
	if lErr != nil {
		if errors.Is(lErr, domain.ErrPasswordIncorrect) || errors.Is(lErr, domain.ErrUserNotFound) {
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
