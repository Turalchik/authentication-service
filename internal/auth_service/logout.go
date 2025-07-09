package auth_service

import (
	"github.com/Turalchik/authentication-service/internal/apperrors"
)

func (authService *AuthService) Logout(accessToken string, userID string) error {
	// заносим access токен в black-list
	if err := authService.tokenRevocationStore.Revoke(accessToken, authService.ttlAccessToken); err != nil {
		return apperrors.ErrCantRevokeToken
	}

	// удаляем refresh токен из базы
	if err := authService.repo.DeleteSessionByUserID(userID); err != nil {
		return apperrors.ErrCantDeleteSession
	}

	return nil
}
