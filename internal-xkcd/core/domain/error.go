package domain

import "errors"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	ErrUserAlreadyAdmin = errors.New("user role already Admin")

	ErrUserRoleUnexpected = errors.New("user role unexpected")

	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")

	ErrLoginIncorrect    = errors.New("login incorrect")
	ErrPasswordIncorrect = errors.New("password incorrect")
)
