package port

import (
	"github.com/golang-jwt/jwt/v5"
	"myapp/internal/core/domain"
)

type TokenService interface {
	CreateToken(user *domain.User) (string, error)
	GetUserByToken(token *jwt.Token) (string, error)
	GetTokenByString(tokenString string) (*jwt.Token, error)
}

type AuthService interface {
	Login(login string, password string) (string, error)
}
