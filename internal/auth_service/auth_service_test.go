package auth_service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/Turalchik/authentication-service/internal/entities/sessions"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) GetSessionByUserID(userID string) (*sessions.Sessions, error) {
	args := m.Called(userID)
	return args.Get(0).(*sessions.Sessions), args.Error(1)
}
func (m *mockRepo) CreateSession(session *sessions.Sessions) error {
	return m.Called(session).Error(0)
}
func (m *mockRepo) DeleteSessionByUserID(userID string) error {
	return m.Called(userID).Error(0)
}
func (m *mockRepo) UpdateRefreshTokenByUserID(userID string, newRefreshTokenHash string) error {
	return m.Called(userID, newRefreshTokenHash).Error(0)
}

type mockTokenRevocationStore struct{ mock.Mock }

func (m *mockTokenRevocationStore) Revoke(tokenID string, ttl time.Duration) error {
	return m.Called(tokenID, ttl).Error(0)
}
func (m *mockTokenRevocationStore) IsRevoked(tokenID string) (bool, error) {
	args := m.Called(tokenID)
	return args.Bool(0), args.Error(1)
}

func TestAuthService_CreateTokens(t *testing.T) {
	repo := new(mockRepo)
	tokenStore := new(mockTokenRevocationStore)
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "", tokenStore)

	t.Run("invalid user id", func(t *testing.T) {
		access, refresh, err := svc.CreateTokens("", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrInvalidUserID)
		assert.Empty(t, access)
		assert.Empty(t, refresh)
	})

	t.Run("user already exists", func(t *testing.T) {
		repo.On("GetSessionByUserID", "u").Return(&sessions.Sessions{}, nil).Once()
		access, refresh, err := svc.CreateTokens("u", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrUserAlreadyExists)
		assert.Empty(t, access)
		assert.Empty(t, refresh)
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo.On("GetSessionByUserID", "u").Return((*sessions.Sessions)(nil), apperrors.ErrUserNotFound).Once()
		repo.On("CreateSession", mock.AnythingOfType("*sessions.Sessions")).Return(nil).Once()
		access, refresh, err := svc.CreateTokens("u", "ua", "ip")
		assert.NoError(t, err)
		assert.NotEmpty(t, access)
		assert.NotEmpty(t, refresh)
		repo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	repo := new(mockRepo)
	tokenStore := new(mockTokenRevocationStore)
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "", tokenStore)

	t.Run("cant revoke token", func(t *testing.T) {
		tokenStore.On("Revoke", "access", time.Minute).Return(errors.New("fail")).Once()
		err := svc.Logout("access", "u")
		assert.ErrorIs(t, err, apperrors.ErrCantRevokeToken)
		tokenStore.AssertExpectations(t)
	})

	t.Run("cant delete session", func(t *testing.T) {
		tokenStore.On("Revoke", "access", time.Minute).Return(nil).Once()
		repo.On("DeleteSessionByUserID", "u").Return(errors.New("fail")).Once()
		err := svc.Logout("access", "u")
		assert.ErrorIs(t, err, apperrors.ErrCantDeleteSession)
		tokenStore.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		tokenStore.On("Revoke", "access", time.Minute).Return(nil).Once()
		repo.On("DeleteSessionByUserID", "u").Return(nil).Once()
		err := svc.Logout("access", "u")
		assert.NoError(t, err)
		tokenStore.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

func TestAuthService_CheckAccessTokenValidity(t *testing.T) {
	repo := new(mockRepo)
	tokenStore := new(mockTokenRevocationStore)
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "", tokenStore)

	t.Run("token revoked", func(t *testing.T) {
		tokenStore.On("IsRevoked", "token").Return(true, nil).Once()
		userID, err := svc.CheckAccessTokenValidity("token")
		assert.ErrorIs(t, err, apperrors.ErrInvalidToken)
		assert.Empty(t, userID)
		tokenStore.AssertExpectations(t)
	})

	t.Run("cant check revocation", func(t *testing.T) {
		tokenStore.On("IsRevoked", "token").Return(false, errors.New("fail")).Once()
		userID, err := svc.CheckAccessTokenValidity("token")
		assert.ErrorIs(t, err, apperrors.ErrCantCheckRevocationToken)
		assert.Empty(t, userID)
		tokenStore.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		tokenStore.On("IsRevoked", "bad").Return(false, nil).Once()
		userID, err := svc.CheckAccessTokenValidity("bad")
		assert.Error(t, err)
		assert.Empty(t, userID)
		tokenStore.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		access, _ := makeJWT("u", time.Minute, []byte("secret"))
		tokenStore.On("IsRevoked", access).Return(false, nil).Once()
		userID, err := svc.CheckAccessTokenValidity(access)
		assert.NoError(t, err)
		assert.Equal(t, "u", userID)
		tokenStore.AssertExpectations(t)
	})
}

func TestAuthService_RefreshTokens(t *testing.T) {
	repo := new(mockRepo)
	tokenStore := new(mockTokenRevocationStore)
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "", tokenStore)
	access, _ := makeJWT("u", time.Minute, []byte("secret"))
	hash, _ := bcrypt.GenerateFromPassword([]byte("refresh"), bcrypt.DefaultCost)
	sess := &sessions.Sessions{UserID: "u", RefreshTokenHash: hash, UserAgent: "ua", IPAddr: "ip"}

	t.Run("token revoked", func(t *testing.T) {
		tokenStore.On("IsRevoked", "revoked").Return(true, nil).Once()
		_, _, err := svc.RefreshTokens("revoked", "refresh", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrInvalidToken)
		tokenStore.AssertExpectations(t)
	})

	t.Run("invalid access token", func(t *testing.T) {
		tokenStore.On("IsRevoked", "bad").Return(false, nil).Once()
		_, _, err := svc.RefreshTokens("bad", "refresh", "ua", "ip")
		assert.Error(t, err)
		tokenStore.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		tokenStore.On("IsRevoked", access).Return(false, nil).Once()
		repo.On("GetSessionByUserID", "u").Return((*sessions.Sessions)(nil), apperrors.ErrUserNotFound).Once()
		_, _, err := svc.RefreshTokens(access, "refresh", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
		tokenStore.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("cant get session", func(t *testing.T) {
		tokenStore.On("IsRevoked", access).Return(false, nil).Once()
		repo.On("GetSessionByUserID", "u").Return((*sessions.Sessions)(nil), errors.New("fail")).Once()
		_, _, err := svc.RefreshTokens(access, "refresh", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrCantGetSession)
		tokenStore.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("tokens dont match", func(t *testing.T) {
		tokenStore.On("IsRevoked", access).Return(false, nil).Once()
		repo.On("GetSessionByUserID", "u").Return(sess, nil).Once()
		_, _, err := svc.RefreshTokens(access, "wrong", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrTokensDontMatch)
		tokenStore.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		tokenStore.On("IsRevoked", access).Return(false, nil).Once()
		repo.On("GetSessionByUserID", "u").Return(sess, nil).Once()
		repo.On("UpdateRefreshTokenByUserID", "u", mock.Anything).Return(nil).Once()
		newAccess, newRefresh, err := svc.RefreshTokens(access, "refresh", "ua", "ip")
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
		assert.NotEmpty(t, newRefresh)
		tokenStore.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}
