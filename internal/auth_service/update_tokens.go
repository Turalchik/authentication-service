package auth_service

import (
	"database/sql"
	"errors"
	"github.com/Turalchik/authentication-service/internal/apperrors"
)

func (authService *AuthService) UpdateTokens(userID string, userAgent string, userIP string) (string, string, error) {
	if userID == "" {
		return "", "", apperrors.ErrInvalidUserID
	}

	session, err := authService.repo.GetSessionByUserID(userID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return "", "", apperrors.ErrUserNotFound
		}
		return "", "", apperrors.ErrCantUpdateTokens
	}

	// TODO
	// подумать над деавторизацией пользователя
	// нужно закидывать access токен в black-лист c ttl = ttlAccessToken

}
