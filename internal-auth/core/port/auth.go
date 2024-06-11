package port

import "myapp/internal-auth/core/domain"

type TokenService interface {
	CreateToken(user *domain.User) (string, error)
	GetUserByTokenString(tokenString string) (string, error)
}

type AuthService interface {
	Login(login string, password string) (string, error)
}
