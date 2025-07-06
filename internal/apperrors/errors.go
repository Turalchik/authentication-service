package apperrors

import (
	"errors"
)

var (
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrUserNotFound      = errors.New("user not found")
	ErrCantCreateTokens  = errors.New("can't create tokens")
	ErrCantCreateSession = errors.New("can't create session")
	ErrCantUpdateTokens  = errors.New("can't update tokens")
	ErrUserAlreadyExists = errors.New("user already exists")
)
