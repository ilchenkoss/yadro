package handlers

import (
	"encoding/json"
	"errors"
	"myapp/internal/core/domain"
	"myapp/internal/core/port"
	"net/http"
)

type UserHandler struct {
	svc port.UserService
	rep port.UserRepository
}

func NewUserHandler(svc port.UserService, rep port.UserRepository) *UserHandler {
	return &UserHandler{
		svc,
		rep,
	}
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

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

	user := domain.User{
		Login:    reqBody.Login,
		Password: reqBody.Password,
	}

	err := uh.svc.Register(&user)
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

func (uh *UserHandler) ToAdmin(w http.ResponseWriter, r *http.Request) {

	type ToAdminRequest struct {
		Login string `json:"login"`
	}

	var reqBody ToAdminRequest
	dErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if dErr != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if len(reqBody.Login) == 0 {
		http.Error(w, "login are required", http.StatusBadRequest)
		return
	}

	taErr := uh.svc.ToAdmin(&domain.User{Login: reqBody.Login})
	if taErr != nil {
		if errors.Is(taErr, domain.ErrUserNotFound) {
			http.Error(w, "login not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "to admin failed", http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode("to admin successful")
	if err != nil {
		//nothing
		return
	}
}
