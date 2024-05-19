package service

import (
	"errors"
	"myapp/internal/core/domain"
	"myapp/internal/core/port"
	"myapp/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

func (us *UserService) Register(user *domain.User) error {

	user.Role = domain.Ordinary

	salt, gsErr := util.GenerateSalt(10)
	if gsErr != nil {
		//domain.ErrLengthMustBePositive
		return gsErr
	}

	user.Salt = salt

	hashedPassword, ghpErr := util.HashPassword(user.Password, salt, user.Role)
	if ghpErr != nil {
		return ghpErr
	}

	user.Password = hashedPassword

	cuErr := us.repo.CreateUser(user)
	if cuErr != nil {
		return cuErr
	}
	return nil
}

func (us *UserService) ToAdmin(user *domain.User) error {

	user, guErr := us.repo.GetUserByLogin(user.Login)
	if guErr != nil {
		// domain.ErrUserNotFound
		return guErr
	}
	if user.Role == domain.Admin {
		return errors.New("user role already Admin")
	}
	user.Role = domain.Admin

	usErr := us.repo.UpdateUser(user)
	if usErr != nil {
		return usErr
	}

	return nil
}

func (us *UserService) RegisterSuperAdmin(user *domain.User) error {

	user.Role = domain.SuperAdmin

	salt, gsErr := util.GenerateSalt(10)
	if gsErr != nil {
		//domain.ErrLengthMustBePositive
		return gsErr
	}

	user.Salt = salt

	hashedPassword, ghpErr := util.HashPassword(user.Password, salt, user.Role)
	if ghpErr != nil {
		return ghpErr
	}
	userPassword := user.Password
	user.Password = hashedPassword

	cuErr := us.repo.CreateUser(user)
	if cuErr != nil {
		if errors.Is(cuErr, domain.ErrUserAlreadyExist) {

			existUser, guErr := us.repo.GetUserByLogin(user.Login)
			if guErr != nil {
				return guErr
			}

			if existUser.Role != domain.SuperAdmin {
				return domain.ErrUserNotSuperAdmin
			}

			cpErr := util.ComparePassword(userPassword, existUser.Salt, existUser.Password)
			if cpErr != nil {
				//domain.ErrPasswordIncorrect
				return cpErr
			}

			return nil
		}
		return cuErr
	}
	return nil
}
