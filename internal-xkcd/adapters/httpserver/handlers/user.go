package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"myapp/internal-xkcd/core/domain"
	"myapp/internal-xkcd/core/port"
	"net/http"
)

type UserHandler struct {
	uc port.UserClient
}

func NewUserHandler(authClient port.UserClient) *UserHandler {
	return &UserHandler{
		uc: authClient,
	}
}

func (uh *UserHandler) ToAdmin(w http.ResponseWriter, r *http.Request) {

	type ToAdminRequest struct {
		ReqUserID int64 `json:"user_id"`
	}

	var reqBody ToAdminRequest
	dErr := json.NewDecoder(r.Body).Decode(&reqBody)
	if dErr != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if crErr := uh.uc.ChangeRole(reqBody.ReqUserID, domain.Admin); crErr != nil {
		fmt.Println(crErr)
		switch {
		case errors.Is(crErr, domain.ErrUserNotFound):
			http.Error(w, "req user not found", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	err := json.NewEncoder(w).Encode("to admin successful")
	if err != nil {
		//nothing
	}
}
