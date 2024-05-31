package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"myapp/internal/core/domain"
	"myapp/internal/core/port/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {

	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	tests := []struct {
		name  string
		mocks func(
			uRepo *mock.MockUserRepository,
			uService *mock.MockUserService,
		)

		requestBody  interface{}
		expectedCode int
	}{
		{
			name:        "Success",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					Register(gomock.Any()).Return(&domain.User{}, nil)
			},
			expectedCode: http.StatusOK,
		}, {
			name:        "Error login and password required",
			requestBody: LoginRequest{Login: "login", Password: ""},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
			},
			expectedCode: http.StatusBadRequest,
		}, {
			name:        "Error user already exist",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					Register(gomock.Any()).Return(&domain.User{}, domain.ErrUserAlreadyExist)
			},
			expectedCode: http.StatusBadRequest,
		}, {
			name:        "Error user already exist",
			requestBody: LoginRequest{Login: "login", Password: "password"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					Register(gomock.Any()).Return(&domain.User{}, errors.New("unhandled error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mock.NewMockUserService(ctrl)
			userRepo := mock.NewMockUserRepository(ctrl)
			tt.mocks(userRepo, userService)

			userHandler := NewUserHandler(userService, userRepo)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(userHandler.Register)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}

func TestUserHandler_ToAdmin(t *testing.T) {

	type LoginRequest struct {
		Login string `json:"login"`
	}

	tests := []struct {
		name  string
		mocks func(
			uRepo *mock.MockUserRepository,
			uService *mock.MockUserService,
		)

		requestBody  interface{}
		expectedCode int
	}{
		{
			name:        "Success",
			requestBody: LoginRequest{Login: "login"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					ToAdmin(gomock.Any()).Return(&domain.User{}, nil)
			},
			expectedCode: http.StatusOK,
		}, {
			name:        "Error user not found",
			requestBody: LoginRequest{Login: "login"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					ToAdmin(gomock.Any()).Return(nil, domain.ErrUserNotFound)
			},
			expectedCode: http.StatusBadRequest,
		}, {
			name:        "Error user already admin",
			requestBody: LoginRequest{Login: "login"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					ToAdmin(gomock.Any()).Return(nil, domain.ErrUserAlreadyAdmin)
			},
			expectedCode: http.StatusBadRequest,
		}, {
			name:        "Error user already admin",
			requestBody: LoginRequest{Login: "login"},
			mocks: func(
				uRepo *mock.MockUserRepository,
				uService *mock.MockUserService,
			) {
				uService.EXPECT().
					ToAdmin(gomock.Any()).Return(nil, errors.New("unhandled error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService := mock.NewMockUserService(ctrl)
			userRepo := mock.NewMockUserRepository(ctrl)
			tt.mocks(userRepo, userService)

			userHandler := NewUserHandler(userService, userRepo)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(userHandler.ToAdmin)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
