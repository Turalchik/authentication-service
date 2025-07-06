package auth_service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func makeJWT(userID string, ttl time.Duration, jwtSecretKey []byte) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return tok.SignedString(jwtSecretKey)
}

func makeTokenInBase64() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(raw)
	return token, nil
}

func claimsFromAccessToken(tokenStr string, jwtSecretKey []byte) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecretKey, nil
	})
	if err != nil || !tok.Valid {
		return nil, apperrors.ErrInvalidToken
	}
	return claims, nil
}

func notifyWebhook(userID string, oldIP string, newIP string, webhookURL string) (resp *http.Response, err error) {
	payload := map[string]string{
		"user_id":     userID,
		"original_ip": oldIP,
		"new_ip":      newIP,
	}
	b, _ := json.Marshal(payload)
	return http.Post(webhookURL, "application/json", io.NopCloser(bytes.NewReader(b)))
}
