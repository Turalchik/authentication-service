package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// мок для AuthService

type mockAuthService struct {
	CreateTokensFunc             func(userID, userAgent, userIP string) (string, string, error)
	RefreshTokensFunc            func(access, refresh, userAgent, userIP string) (string, string, error)
	LogoutFunc                   func(access, userID string) error
	CheckAccessTokenValidityFunc func(token string) (string, error)
}

func (m *mockAuthService) CreateTokens(userID, userAgent, userIP string) (string, string, error) {
	return m.CreateTokensFunc(userID, userAgent, userIP)
}
func (m *mockAuthService) RefreshTokens(access, refresh, userAgent, userIP string) (string, string, error) {
	if m.RefreshTokensFunc != nil {
		return m.RefreshTokensFunc(access, refresh, userAgent, userIP)
	}
	return "", "", nil
}
func (m *mockAuthService) Logout(access, userID string) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(access, userID)
	}
	return nil
}
func (m *mockAuthService) CheckAccessTokenValidity(token string) (string, error) {
	if m.CheckAccessTokenValidityFunc != nil {
		return m.CheckAccessTokenValidityFunc(token)
	}
	return "", nil
}

func TestHttpHandler_CreateTokens(t *testing.T) {
	handler := &HttpHandler{
		authService: &mockAuthService{
			CreateTokensFunc: func(userID, userAgent, userIP string) (string, string, error) {
				if userID == "fail" {
					return "", "", errors.New("fail")
				}
				return "access", "refresh", nil
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/tokens?user_id=123", nil)
		req.Header.Set("User-Agent", "test-agent")
		req.RemoteAddr = "1.2.3.4:5678"
		rw := httptest.NewRecorder()

		handler.CreateTokens(rw, req)

		assert.Equal(t, http.StatusOK, rw.Code)
		var resp map[string]string
		err := json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "access", resp["access_token"])
		assert.Equal(t, "refresh", resp["refresh_token"])
	})

	t.Run("missing user_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/tokens", nil)
		rw := httptest.NewRecorder()
		handler.CreateTokens(rw, req)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "user_id required")
	})

	t.Run("service error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/tokens?user_id=fail", nil)
		rw := httptest.NewRecorder()
		handler.CreateTokens(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func TestHttpHandler_RefreshTokens(t *testing.T) {
	handler := &HttpHandler{
		authService: &mockAuthService{
			RefreshTokensFunc: func(access, refresh, userAgent, userIP string) (string, string, error) {
				if access == "bad" {
					return "", "", errors.New("bad token")
				}
				return "new_access", "new_refresh", nil
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		body := `{"access_token":"a","refresh_token":"r"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(body))
		req.Header.Set("User-Agent", "test-agent")
		req.RemoteAddr = "1.2.3.4:5678"
		rw := httptest.NewRecorder()
		handler.RefreshTokens(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		var resp map[string]string
		err := json.Unmarshal(rw.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "new_access", resp["access_token"])
		assert.Equal(t, "new_refresh", resp["refresh_token"])
	})

	t.Run("bad body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader("{"))
		rw := httptest.NewRecorder()
		handler.RefreshTokens(rw, req)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "Invalid request body")
	})

	t.Run("service error", func(t *testing.T) {
		body := `{"access_token":"bad","refresh_token":"r"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(body))
		rw := httptest.NewRecorder()
		handler.RefreshTokens(rw, req)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "Invalid request")
	})
}

func TestHttpHandler_Logout(t *testing.T) {
	handler := &HttpHandler{
		authService: &mockAuthService{
			LogoutFunc: func(access, userID string) error {
				if access == "bad" {
					return errors.New("bad token")
				}
				return nil
			},
		},
	}

	ctx := context.WithValue(context.Background(), "args", map[string]string{"accessToken": "good", "userID": "u"})
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil).WithContext(ctx)
		rw := httptest.NewRecorder()
		handler.Logout(rw, req)
		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	t.Run("service error", func(t *testing.T) {
		ctxBad := context.WithValue(context.Background(), "args", map[string]string{"accessToken": "bad", "userID": "u"})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil).WithContext(ctxBad)
		rw := httptest.NewRecorder()
		handler.Logout(rw, req)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "invalid access token")
	})
}

func TestHttpHandler_Guid(t *testing.T) {
	handler := &HttpHandler{}
	ctx := context.WithValue(context.Background(), "args", map[string]string{"userID": "u"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/guid", nil).WithContext(ctx)
	rw := httptest.NewRecorder()
	handler.Guid(rw, req)
	assert.Equal(t, http.StatusOK, rw.Code)
	var resp map[string]string
	_ = json.Unmarshal(rw.Body.Bytes(), &resp)
	assert.Equal(t, "u", resp["user_id"])
}

func TestHttpHandler_AuthMiddleware(t *testing.T) {
	handler := &HttpHandler{
		authService: &mockAuthService{
			CheckAccessTokenValidityFunc: func(token string) (string, error) {
				if token == "bad" {
					return "", errors.New("bad token")
				}
				return "user", nil
			},
		},
	}

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rw := httptest.NewRecorder()
		handler.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("should not call next")
		})).ServeHTTP(rw, req)
		assert.Equal(t, http.StatusUnauthorized, rw.Code)
		assert.Contains(t, rw.Body.String(), "missing token")
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer bad")
		rw := httptest.NewRecorder()
		handler.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("should not call next")
		})).ServeHTTP(rw, req)
		assert.Equal(t, http.StatusUnauthorized, rw.Code)
		assert.Contains(t, rw.Body.String(), "invalid token")
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer good")
		rw := httptest.NewRecorder()
		called := false
		handler.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			args := r.Context().Value("args").(map[string]string)
			assert.Equal(t, "user", args["userID"])
			assert.Equal(t, "good", args["accessToken"])
		})).ServeHTTP(rw, req)
		assert.True(t, called)
	})
}
