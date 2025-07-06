package handlers

type AuthService interface {
	CreateTokens(userID string) (string, string, error)
}
