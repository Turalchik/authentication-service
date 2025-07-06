package auth_service

import "time"

type AuthService struct {
	repo Repo

	ttlAccessToken time.Duration

	jwtSecretKey []byte
	webhookURL   string
}

func NewAuthService(repo Repo, ttlAccessToken time.Duration, jwtSecretKey []byte, webhookURL string) *AuthService {
	return &AuthService{
		repo:           repo,
		ttlAccessToken: ttlAccessToken,
		jwtSecretKey:   jwtSecretKey,
		webhookURL:     webhookURL,
	}
}
