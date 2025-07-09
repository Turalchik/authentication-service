package token_revocation_store

import (
	"context"
	"time"
)

// Revoke — устанавливает в Redis ключ <prefix><tokenID> = true с TTL
func (revocationStore *TokenRevocationStore) Revoke(token string, ttl time.Duration) error {
	key := revocationStore.keyPrefix + token
	return revocationStore.client.Set(context.Background(), key, "1", ttl).Err()
}
