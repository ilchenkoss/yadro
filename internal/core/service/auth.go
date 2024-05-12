package service

import (
	"myapp/internal/core/port"
	"myapp/internal/core/util"
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

func NewAuthService(repo port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		repo,
		ts,
	}
}

func (as *AuthService) Login(login string, password string) (string, error) {

	user, guErr := as.repo.GetUserByLogin(login)
	if guErr != nil {
		//domain.ErrUserNotFound
		return "", guErr
	}

	cpErr := util.ComparePassword(password, user.Salt, user.Password)
	if cpErr != nil {
		//domain.ErrPasswordIncorrect
		return "", cpErr
	}

	accessToken, ctErr := as.ts.CreateToken(user)
	if ctErr != nil {
		return "", ctErr
	}

	return accessToken, nil
}
