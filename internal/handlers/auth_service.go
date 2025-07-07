package handlers

type AuthService interface {
	CreateTokens(userID string, userAgent string, userIP string) (string, string, error)
	RefreshTokens(accessToken string, refreshToken string, userAgent string, userIP string) (string, string, error)
	Logout(accessToken string, userID string) error
	CheckAccessTokenValidity(accessToken string) (string, error)
}
