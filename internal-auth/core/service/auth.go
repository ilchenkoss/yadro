package service

import (
	"myapp/internal-auth/core/domain"
	"myapp/internal-auth/core/port"
	"myapp/internal-auth/core/service/util"
)

type AuthService struct {
	ur port.UserRepository
	ts port.TokenService
}

func NewAuthService(ur port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		ur: ur,
		ts: ts,
	}
}

func (as *AuthService) Login(login string, password string) (string, error) {

	user, guErr := as.ur.GetUserByLogin(login)
	if guErr != nil {
		//domain.ErrUserNotFound
		return "", guErr
	}

	cpErr := util.ComparePassword(password, user.Salt, user.Password)
	if cpErr != nil {
		//domain.ErrPasswordIncorrect
		return "", cpErr
	}

	accessToken, ctErr := as.ts.CreateToken(user.ID)
	if ctErr != nil {
		return "", ctErr
	}

	return accessToken, nil
}

func (as *AuthService) Register(login string, password string, role domain.UserRole) (int64, error) {

	user := domain.User{}

	user.Login = login
	user.Role = role

	salt, gsErr := util.GenerateSalt(10)
	if gsErr != nil {
		//domain.ErrLengthMustBePositive
		return 0, gsErr
	}

	user.Salt = salt

	hashedPassword, ghpErr := util.HashPassword(password, salt, user.Role)
	if ghpErr != nil {
		return 0, ghpErr
	}

	user.Password = hashedPassword

	uUser, cuErr := as.ur.CreateUser(&user)
	if cuErr != nil {
		//domain.ErrUserAlreadyExist
		return 0, cuErr
	}

	return uUser.ID, nil
}

func (us *AuthService) UserID(token string) (int64, error) {
	userID, guIDErr := us.ts.GetUserID(token)
	if guIDErr != nil {
		//domain.ErrTokenNotValid
		//domain.ErrTokenExpired
		return 0, guIDErr
	}
	return userID, nil
}
