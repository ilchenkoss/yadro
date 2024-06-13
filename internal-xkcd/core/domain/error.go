package domain

import "errors"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	ErrTokenExpired  = errors.New("token expired")
	ErrTokenNotValid = errors.New("token not valid")

	ErrUserAlreadyAdmin = errors.New("user role already Admin")

	ErrUserRoleUnexpected = errors.New("user role unexpected")

	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")

	ErrLoginIncorrect    = errors.New("login incorrect")
	ErrPasswordIncorrect = errors.New("password incorrect")
)
