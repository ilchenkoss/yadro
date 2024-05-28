package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"myapp/internal-api/core/domain"
	"myapp/internal-api/core/port/mock"
	"myapp/internal-api/core/util"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	userLogin := "user1"
	userPassword := "strong_password"

	hashedPassword, _ := util.HashPassword(userPassword, "salt", domain.Ordinary)

	user := &domain.User{Login: userLogin, Password: hashedPassword, Salt: "salt"}

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
			tokenService *mock.MockTokenService,
		)

		userLoginInput    string
		userPasswordInput string

		expectToken string
		expectErr   error
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
				tokenService *mock.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(user, nil)
				tokenService.EXPECT().
					CreateToken(gomock.Any()).Return("token", nil)
			},

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectToken: "token",
			expectErr:   nil,
		}, {
			desc: "Error user not found",
			mocks: func(
				userRepo *mock.MockUserRepository,
				tokenService *mock.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(nil, domain.ErrUserNotFound)
			},

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectToken: "",
			expectErr:   domain.ErrUserNotFound,
		}, {
			desc: "Error password incorrect",
			mocks: func(
				userRepo *mock.MockUserRepository,
				tokenService *mock.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(user, nil)
			},

			userLoginInput:    userLogin,
			userPasswordInput: "bad password",

			expectToken: "",
			expectErr:   domain.ErrPasswordIncorrect,
		}, {
			desc: "Error token service",
			mocks: func(
				userRepo *mock.MockUserRepository,
				tokenService *mock.MockTokenService,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(user, nil)
				tokenService.EXPECT().
					CreateToken(gomock.Any()).Return("", errors.New("token error"))
			},

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectToken: "",
			expectErr:   errors.New("token error"),
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)
			tokenService := mock.NewMockTokenService(ctrl)

			tc.mocks(userRepo, tokenService)

			userService := NewAuthService(userRepo, tokenService)

			token, err := userService.Login(tc.userLoginInput, tc.userPasswordInput)

			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, tc.expectToken, token)
		})
	}
}
