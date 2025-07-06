package auth_service

import "time"

type AuthService struct {
	repo Repo

	ttlAccessToken time.Duration

	jwtSecretKey []byte
	webHookURL   string
}
