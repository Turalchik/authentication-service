package auth_service

import "github.com/Turalchik/authentication-service/internal/entities/sessions"

type Repo interface {
	GetSessionByUserID(userID string) (*sessions.Sessions, error)
	CreateSession(session *sessions.Sessions) error
	DeleteSessionByUserID(userID string) error
	UpdateRefreshTokenByUserID(userID string, newRefreshTokenHash string) error
}
