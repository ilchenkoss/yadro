package domain

import "errors"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	ErrTokenExpired = errors.New("token expired")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotSuperAdmin = errors.New("user not super admin")

	ErrPasswordIncorrect = errors.New("password incorrect")
)
