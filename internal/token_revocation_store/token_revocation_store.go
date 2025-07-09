package token_revocation_store

import (
	"github.com/go-redis/redis/v8"
)

type TokenRevocationStore struct {
	client    *redis.Client
	keyPrefix string
}

func NewTokenRevocationStore(addr, password string, db int, keyPrefix string) *TokenRevocationStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &TokenRevocationStore{
		client:    rdb,
		keyPrefix: keyPrefix,
	}
}
