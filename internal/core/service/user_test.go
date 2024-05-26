package service

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"myapp/internal/core/domain"
	"myapp/internal/core/port/mock"
	"myapp/internal/core/util"
	"strings"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	userLogin := "user1"
	userPassword := "strong_password"

	userInput := domain.User{
		Login:    userLogin,
		Password: userPassword,
	}

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
		)
		userInput         domain.User
		userLoginInput    string
		userPasswordInput string

		expectErr error
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any()).Return(nil)
			},
			userInput: userInput,

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectErr: nil,
		},
		{
			desc: "Err user exist",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any()).Return(domain.ErrUserAlreadyExist)
			},
			userInput:         userInput,
			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectErr: domain.ErrUserAlreadyExist,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)

			tc.mocks(userRepo)

			userService := NewUserService(userRepo)

			user, err := userService.Register(&tc.userInput)

			assert.Equal(t, tc.expectErr, err)

			if strings.Contains(tc.desc, "Success") {
				assert.NoError(t, util.ComparePassword(tc.userPasswordInput, user.Salt, user.Password))
				assert.NotNil(t, user.Salt)
				assert.Equal(t, user.Role, domain.Ordinary)
			}

		})
	}
}

func TestUserService_ToAdmin(t *testing.T) {
	userLogin := "user1"

	userInput := domain.User{
		Login: userLogin,
	}

	userExpect := domain.User{
		Login: userInput.Login,
		Role:  domain.Admin,
	}

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
		)
		userInput  domain.User
		userExpect domain.User

		expectErr error
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(&userInput, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any()).Return(nil)
			},
			userInput:  userInput,
			userExpect: userExpect,

			expectErr: nil,
		}, {
			desc: "Error user already admin",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(&userInput, domain.ErrUserAlreadyAdmin)
			},
			userInput:  userInput,
			userExpect: userExpect,

			expectErr: domain.ErrUserAlreadyAdmin,
		}, {
			desc: "Error user not found",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(&userInput, domain.ErrUserNotFound)
			},
			userInput:  userInput,
			userExpect: userExpect,

			expectErr: domain.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)

			tc.mocks(userRepo)

			userService := NewUserService(userRepo)

			user, err := userService.ToAdmin(&tc.userInput)

			assert.Equal(t, tc.expectErr, err)
			if strings.Contains(tc.desc, "Success") {
				assert.Equal(t, &tc.userExpect, user)
			}
		})
	}
}

func TestUserService_RegisterSuperAdmin(t *testing.T) {
	userLogin := "user1"
	userPassword := "strong_password"

	userInput := domain.User{
		Login:    userLogin,
		Password: userPassword,
	}

	hashedPassword, _ := util.HashPassword(userPassword, "salt", domain.Ordinary)

	userExistIsSuperAdmin := domain.User{
		Login:    userLogin,
		Salt:     "salt",
		Password: hashedPassword,
		Role:     domain.SuperAdmin,
	}

	userExistIsNotAdmin := domain.User{
		Login:    userLogin,
		Password: hashedPassword,
		Role:     domain.Ordinary,
	}

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
		)
		userInput         domain.User
		userLoginInput    string
		userPasswordInput string

		expectRole domain.UserRole
		expectErr  error
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any()).Return(nil)
			},
			userInput: userInput,

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectRole: domain.SuperAdmin,
			expectErr:  nil,
		}, {
			desc: "Error user already exist and super admin",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any()).Return(domain.ErrUserAlreadyExist)
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(&userExistIsSuperAdmin, nil)
			},
			userInput: userInput,

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectErr: domain.ErrUserAlreadyExist,
		}, {
			desc: "Error user already exist and not super admin",
			mocks: func(
				userRepo *mock.MockUserRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any()).Return(domain.ErrUserAlreadyExist)
				userRepo.EXPECT().
					GetUserByLogin(gomock.Any()).Return(&userExistIsNotAdmin, nil)
			},
			userInput: userInput,

			userLoginInput:    userLogin,
			userPasswordInput: userPassword,

			expectErr: domain.ErrUserNotSuperAdmin,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.desc, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)

			tc.mocks(userRepo)

			userService := NewUserService(userRepo)

			user, err := userService.RegisterSuperAdmin(&tc.userInput)

			assert.Equal(t, tc.expectErr, err)

			if strings.Contains(tc.desc, "Success") {
				assert.NoError(t, util.ComparePassword(tc.userPasswordInput, user.Salt, user.Password))
				assert.Equal(t, tc.expectRole, user.Role)
				assert.NotNil(t, user.Salt)
			}

		})
	}
}
