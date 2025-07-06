package handlers

type AuthService interface {
	CreateTokens(userID string, userAgent string, userIP string) (string, string, error)
	UpdateTokens(accessToken string, refreshToken string, userAgent string, userIP string) (string, string, error)
}
