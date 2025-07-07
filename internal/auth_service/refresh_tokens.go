package auth_service

import (
	"database/sql"
	"errors"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (authService *AuthService) RefreshTokens(accessToken string, refreshToken string, userAgent string, ipAddr string) (string, string, error) {
	// извлечь claims из access token
	claims, err := claimsFromAccessToken(accessToken, authService.jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// найти соответствующий refresh токен
	session, err := authService.repo.GetSessionByUserID(claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", apperrors.ErrUserNotFound
		}
		return "", "", apperrors.ErrCantGetSession
	}

	// проверяем на соответствие refresh токены
	if bcrypt.CompareHashAndPassword(session.RefreshTokenHash, []byte(refreshToken)) != nil {
		return "", "", apperrors.ErrTokensDontMatch
	}

	// проверить userAgent
	if userAgent != session.UserAgent {
		if err = authService.repo.DeleteSessionByUserID(claims.UserID); err != nil {
			return "", "", apperrors.ErrCantDeleteSession
		}
		// TODO
		// подумать над деавторизацией пользователя
		// нужно закидывать access токен в black-лист (redis) c ttl = ttlAccessToken
	}

	// проверить userIP
	if ipAddr != session.IPAddr {
		go func() {
			_, err := notifyWebhook(claims.UserID, session.IPAddr, ipAddr, authService.webhookURL)
			if err != nil {
				log.Printf("can't notify webhook with error: %s\n", err.Error())
			}
		}()
	}

	// TODO
	// тут тоже нужно старый access токен занести в black-list
	newAccessToken, err := makeJWT(claims.UserID, authService.ttlAccessToken, authService.jwtSecretKey)
	if err != nil {
		return "", "", apperrors.ErrCantCreateTokens
	}

	newRefreshToken, err := makeTokenInBase64()
	if err != nil {
		return "", "", apperrors.ErrCantCreateTokens
	}
	newRefreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", apperrors.ErrCantCreateTokens
	}

	if err = authService.repo.UpdateRefreshTokenByUserID(claims.UserID, newRefreshTokenHash); err != nil {
		return "", "", apperrors.ErrCantUpdateTokens
	}

	return newAccessToken, newRefreshToken, nil
}
