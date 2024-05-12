package domain

import "errors"

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")

	ErrTokenExpired = errors.New("token expired")
	ErrUserNotFound = errors.New("user not found")

	ErrPasswordIncorrect = errors.New("password incorrect")
)
