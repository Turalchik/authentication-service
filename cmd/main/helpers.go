package main

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	TTLAccessToken time.Duration
	JWTSecretKey   []byte
	WebhookURL     string
}

func GetConfigFromEnv() (*Config, error) {
	ttlAccessTokenStr := os.Getenv("TTL_ACCESS_TOKEN")
	ttlAccessToken, err := strconv.Atoi(ttlAccessTokenStr)
	if err != nil {
		return nil, err
	}

	jwtSecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	webhookURL := os.Getenv("WEBHOOK_URL")

	cfg := &Config{
		TTLAccessToken: time.Second * time.Duration(ttlAccessToken),
		JWTSecretKey:   jwtSecretKey,
		WebhookURL:     webhookURL,
	}

	return cfg, nil
}
