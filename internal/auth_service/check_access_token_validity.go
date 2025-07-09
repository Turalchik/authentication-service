package auth_service

import "github.com/Turalchik/authentication-service/internal/apperrors"

func (authService *AuthService) CheckAccessTokenValidity(accessToken string) (string, error) {
	isRevoked, err := authService.tokenRevocationStore.IsRevoked(accessToken)
	if err != nil {
		return "", apperrors.ErrCantCheckRevocationToken
	}
	if isRevoked {
		return "", apperrors.ErrInvalidToken
	}

	claims, err := claimsFromAccessToken(accessToken, authService.jwtSecretKey)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
