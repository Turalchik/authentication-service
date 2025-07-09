package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type tokensResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type guidResp struct {
	UserID string `json:"user_id"`
}

func waitForAPI(t *testing.T, url string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 404 { // 404 значит сервер поднялся, но ручка не найдена
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("API not available at %s after %s", url, timeout)
}

func Test_AuthService_HappyPath(t *testing.T) {
	// Ждём, пока API реально поднимется
	waitForAPI(t, "http://localhost:8080/api/v1/unknown", 30*time.Second)

	client := &http.Client{}
	userID := "123e4567-e89b-12d3-a456-426614174000"

	// 1. Получить токены
	resp, err := client.Get("http://localhost:8080/api/v1/auth/tokens?user_id=" + userID)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var tokens tokensResp
	require.NoError(t, json.Unmarshal(body, &tokens))
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)

	// 2. Обновить токены
	refreshBody, _ := json.Marshal(map[string]string{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
	resp, err = client.Post("http://localhost:8080/api/v1/auth/tokens/refresh", "application/json", bytes.NewReader(refreshBody))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var tokens2 tokensResp
	require.NoError(t, json.Unmarshal(body, &tokens2))
	require.NotEmpty(t, tokens2.AccessToken)
	require.NotEmpty(t, tokens2.RefreshToken)

	// 3. Получить GUID
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/auth/guid", nil)
	req.Header.Set("Authorization", "Bearer "+tokens2.AccessToken)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var guid guidResp
	require.NoError(t, json.Unmarshal(body, &guid))
	require.Equal(t, userID, guid.UserID)

	// 4. Logout
	req, _ = http.NewRequest("POST", "http://localhost:8080/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+tokens2.AccessToken)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, 204, resp.StatusCode)
	resp.Body.Close()

	// 5. Попытка refresh после logout (должна быть ошибка)
	refreshBody, _ = json.Marshal(map[string]string{
		"access_token":  tokens2.AccessToken,
		"refresh_token": tokens2.RefreshToken,
	})
	resp, err = client.Post("http://localhost:8080/api/v1/auth/tokens/refresh", "application/json", bytes.NewReader(refreshBody))
	require.NoError(t, err)
	require.NotEqual(t, 200, resp.StatusCode)
	resp.Body.Close()
}
