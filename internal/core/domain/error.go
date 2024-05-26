package domain

import "errors"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	ErrLengthMustBePositive = errors.New("length must be a positive")

	ErrTokenExpired  = errors.New("token expired")
	ErrTokenNotValid = errors.New("token not valid")

	ErrUserAlreadyAdmin = errors.New("user role already Admin")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotSuperAdmin = errors.New("user not super admin")

	ErrPasswordIncorrect = errors.New("password incorrect")
)
