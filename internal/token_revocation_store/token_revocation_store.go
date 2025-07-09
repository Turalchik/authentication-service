package token_revocation_store

import (
	"github.com/go-redis/redis/v8"
)

type TokenRevocationStore struct {
	client    *redis.Client
	keyPrefix string
}

func NewTokenRevocationStore(client *redis.Client, keyPrefix string) *TokenRevocationStore {
	return &TokenRevocationStore{
		client:    client,
		keyPrefix: keyPrefix,
	}
}
