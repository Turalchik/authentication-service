package auth_service

import (
	"database/sql"
	"errors"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/Turalchik/authentication-service/internal/entities/sessions"
	"golang.org/x/crypto/bcrypt"
)

func (authService *AuthService) CreateTokens(userID string, userAgent string, userIP string) (string, string, error) {
	if userID == "" {
		return "", "", apperrors.ErrInvalidUserID
	}

	_, err := authService.repo.GetSessionByUserID(userID)
	if err == nil {
		return "", "", apperrors.ErrUserAlreadyExists
	}

	// case when the user does not have a token yet
	if errors.Is(sql.ErrNoRows, err) {
		// создаем токены (access и refresh)
		accessToken, err := makeJWT(userID, authService.ttlAccessToken, authService.jwtSecretKey)
		if err != nil {
			return "", "", apperrors.ErrCantCreateTokens
		}

		// создаём refresh токен
		refreshToken, err := makeTokenInBase64()
		if err != nil {
			return "", "", apperrors.ErrCantCreateTokens
		}

		// хэшируем refresh токен
		refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
		if err != nil {
			return "", "", apperrors.ErrCantCreateTokens
		}

		// создаём сессию и сохраняем её в базу
		newSession := &sessions.Sessions{
			UserID:           userID,
			RefreshTokenHash: refreshTokenHash,
			UserAgent:        userAgent,
			IssuedIP:         userIP,
		}
		err = authService.repo.CreateSession(newSession)
		if err != nil {
			return "", "", apperrors.ErrCantCreateSession
		}

		// возвращаем токены
		return accessToken, refreshToken, nil
	}

	return "", "", err
}
