package domain

import "errors"

var (
	ErrInvalidID    = errors.New("invalid id")
	ErrUserNotFound = errors.New("user not found")
)
