package service

import (
	"myapp/internal-auth/core/domain"
	"myapp/internal-auth/core/port"
)

type UserService struct {
	ur port.UserRepository
	as port.AuthService
}

func NewUserService(ur port.UserRepository, as port.AuthService) *UserService {
	return &UserService{
		ur: ur,
		as: as,
	}
}

func (us *UserService) UserRole(userID int64) (domain.UserRole, error) {
	user, guErr := us.ur.GetUserByUserID(userID)
	if guErr != nil {
		//domain.ErrUserNotFound
		return "", guErr
	}
	return user.Role, nil
}

func (us *UserService) UpdateUserRole(reqUserID int64, reqRole domain.UserRole) error {

	reqUser, guErr := us.ur.GetUserByUserID(reqUserID)
	if guErr != nil {
		//domain.ErrUserNotFound
		return guErr
	}

	reqUser.Role = reqRole

	if uuErr := us.ur.UpdateUserByUserID(reqUser); uuErr != nil {
		return uuErr
	}
	return nil
}
