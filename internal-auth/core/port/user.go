package port

import (
	"myapp/internal-auth/core/domain"
)

type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByLogin(login string) (*domain.User, error)
	GetUserByUserID(userID int64) (*domain.User, error)
	UpdateUserByUserID(user *domain.User) error
}

type UserService interface {
	UserRole(userID int64) (domain.UserRole, error)
	UpdateUserRole(reqUserID int64, reqRole domain.UserRole) error
}
