package auth_service

import "time"

type TokenRevocationStore interface {
	Revoke(token string, ttl time.Duration) error
	IsRevoked(token string) (bool, error)
}
