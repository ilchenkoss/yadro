package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"myapp/internal-api/config"
	"myapp/internal-api/core/domain"
	"time"
)

type TokenService struct {
	SecretKey []byte
	Duration  time.Duration
	Method    jwt.SigningMethod
}

func NewTokenService(cfg config.HttpServerConfig) *TokenService {

	duration := time.Duration(cfg.TokenMaxTime) * time.Minute
	//load secret key
	secretKey := []byte("SuperSecretKey")
	method := jwt.SigningMethodHS256

	return &TokenService{
		secretKey,
		duration,
		method,
	}
}

func (ts *TokenService) CreateToken(user *domain.User) (string, error) {

	token := jwt.NewWithClaims(ts.Method, jwt.RegisteredClaims{
		Subject: user.Login,
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(ts.Duration),
		},
	})

	tokenString, err := token.SignedString(ts.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (ts *TokenService) GetUserByTokenString(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return ts.SecretKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return "", domain.ErrTokenNotValid
		case errors.Is(err, jwt.ErrTokenExpired):
			return "", domain.ErrTokenExpired
		default:
			return "", err
		}
	}

	if token == nil {
		return "", domain.ErrTokenNotValid
	}

	return token.Claims.GetSubject()
}
