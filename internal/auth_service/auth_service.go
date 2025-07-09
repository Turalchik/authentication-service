package auth_service

import "time"

type AuthService struct {
	repo                 Repo
	tokenRevocationStore TokenRevocationStore

	ttlAccessToken time.Duration

	jwtSecretKey []byte
	webhookURL   string
}

func NewAuthService(

	repo Repo,
	ttlAccessToken time.Duration,
	jwtSecretKey []byte,
	webhookURL string,
	tokenRevocationStore TokenRevocationStore,

) *AuthService {

	return &AuthService{
		repo:                 repo,
		tokenRevocationStore: tokenRevocationStore,
		ttlAccessToken:       ttlAccessToken,
		jwtSecretKey:         jwtSecretKey,
		webhookURL:           webhookURL,
	}
}
