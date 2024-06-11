package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"myapp/internal-auth/config"
	"myapp/internal-auth/core/domain"
	"strconv"
	"time"
)

type TokenService struct {
	SecretKey []byte
	Duration  time.Duration
	Method    jwt.SigningMethod
}

func NewTokenService(cfg config.Config) *TokenService {

	duration := cfg.TokenTTL
	//load secret key
	secretKey := []byte("SuperSecretKey")
	method := jwt.SigningMethodHS256

	return &TokenService{
		secretKey,
		duration,
		method,
	}
}

func (ts *TokenService) CreateToken(userID int64) (string, error) {

	userIDString := strconv.FormatInt(userID, 10)

	token := jwt.NewWithClaims(ts.Method, jwt.RegisteredClaims{
		Subject: userIDString,
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

func (ts *TokenService) GetUserID(token string) (int64, error) {

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return ts.SecretKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return 0, domain.ErrTokenNotValid
		case errors.Is(err, jwt.ErrTokenExpired):
			return 0, domain.ErrTokenExpired
		default:
			return 0, err
		}
	}

	if t == nil {
		return 0, domain.ErrTokenNotValid
	}

	userIDString, gErr := t.Claims.GetSubject()
	if gErr != nil {
		return 0, gErr
	}

	userID, pErr := strconv.ParseInt(userIDString, 10, 64)
	if pErr != nil {
		return 0, pErr
	}

	return userID, nil
}
