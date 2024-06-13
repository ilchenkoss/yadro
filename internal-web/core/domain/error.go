package domain

import "errors"

var (
	ErrUnauthorized   = errors.New("")
	ErrToManyRequests = errors.New("")
	ErrAuthFailed     = errors.New("")
)
