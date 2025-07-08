package auth_service

func (authService *AuthService) CheckAccessTokenValidity(accessToken string) (string, error) {
	claims, err := claimsFromAccessToken(accessToken, authService.jwtSecretKey)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
