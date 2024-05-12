package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"myapp/internal/config"
	"myapp/internal/core/domain"
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
			time.Now().Add(ts.Duration),
		},
	})

	tokenString, err := token.SignedString(ts.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (ts *TokenService) GetUserByToken(token *jwt.Token) (string, error) {
	return token.Claims.GetSubject()
}

func (ts *TokenService) GetTokenByString(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return ts.SecretKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, domain.ErrTokenNotValid
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, domain.ErrTokenExpired
		default:
			return nil, err
		}
	}

	if token == nil {
		return nil, domain.ErrTokenNotValid
	}

	return token, err
}
