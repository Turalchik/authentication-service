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

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func GetConfigFromEnv() (*Config, error) {
	ttlAccessToken, err := strconv.Atoi(os.Getenv("TTL_ACCESS_TOKEN"))
	if err != nil {
		return nil, err
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		TTLAccessToken: time.Second * time.Duration(ttlAccessToken),
		JWTSecretKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		WebhookURL:     os.Getenv("WEBHOOK_URL"),
		
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
	}

	return cfg, nil
}
