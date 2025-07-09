package apperrors

import (
	"errors"
)

var (
	ErrInvalidUserID            = errors.New("invalid user id")
	ErrUserNotFound             = errors.New("user not found")
	ErrCantCreateTokens         = errors.New("can't create tokens")
	ErrCantCreateSession        = errors.New("can't create session")
	ErrCantUpdateTokens         = errors.New("can't update tokens")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrInvalidToken             = errors.New("invalid token")
	ErrCantGetSession           = errors.New("can't get session")
	ErrCantDeleteSession        = errors.New("can't delete session")
	ErrTokensDontMatch          = errors.New("tokens don't match")
	ErrCantBuildSQLQuery        = errors.New("cant build sql query")
	ErrCantExecSQLQuery         = errors.New("can't exec sql query")
	ErrCantOpenDatabase         = errors.New("can't open database")
	ErrCantCheckRevocationToken = errors.New("can't verify the revocation of the token")
	ErrCantRevokeToken          = errors.New("can't revoke token")
)
