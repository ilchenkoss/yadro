package port

import (
	"myapp/internal/core/domain"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByLogin(login string) (*domain.User, error)
	UpdateUser(user *domain.User) error
}

type UserService interface {
	Register(user *domain.User) (*domain.User, error)
	ToAdmin(user *domain.User) (*domain.User, error)
	RegisterSuperAdmin(user *domain.User) (*domain.User, error)
}
