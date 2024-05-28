package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"myapp/internal-api/core/domain"
	"myapp/internal-api/core/port/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler(t *testing.T) {

	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	tests := []struct {
		name  string
		mocks func(
			aService *mock.MockAuthService,
		)

		requestBody  interface{}
		expectedCode int
	}{
		{
			name:        "Success",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				aService *mock.MockAuthService,
			) {
				aService.EXPECT().
					Login(gomock.Any(), gomock.Any()).Return("validToken", nil)
			},
			expectedCode: http.StatusOK,
		}, {
			name:        "Error decode request",
			requestBody: "bad request body",
			mocks: func(
				aService *mock.MockAuthService,
			) {
			},
			expectedCode: http.StatusBadRequest,
		}, {
			name:        "Error password len == 0",
			requestBody: LoginRequest{Login: "login", Password: ""},
			mocks: func(
				aService *mock.MockAuthService,
			) {
			},
			expectedCode: http.StatusBadRequest,
		}, {
			name:        "Error password incorrect",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				aService *mock.MockAuthService,
			) {
				aService.EXPECT().
					Login(gomock.Any(), gomock.Any()).Return("", domain.ErrPasswordIncorrect)
			},
			expectedCode: http.StatusUnauthorized,
		}, {
			name:        "Error user not found",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				aService *mock.MockAuthService,
			) {
				aService.EXPECT().
					Login(gomock.Any(), gomock.Any()).Return("", domain.ErrUserNotFound)
			},
			expectedCode: http.StatusUnauthorized,
		}, {
			name:        "Error unhandled",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				aService *mock.MockAuthService,
			) {
				aService.EXPECT().
					Login(gomock.Any(), gomock.Any()).Return("", errors.New("new error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := mock.NewMockAuthService(ctrl)
			tt.mocks(authService)

			authHandler := NewAuthHandler(authService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(authHandler.Login)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
