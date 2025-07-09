package auth_service

import "time"

type TokenRevocationStore interface {
	Revoke(tokenID string, ttl time.Duration) error
	IsRevoked(tokenID string) (bool, error)
}
