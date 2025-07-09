package token_revocation_store

import "context"

// IsRevoked — проверяет наличие ключа в Redis
func (revocationStore *TokenRevocationStore) IsRevoked(token string) (bool, error) {
	key := revocationStore.keyPrefix + token
	exists, err := revocationStore.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
