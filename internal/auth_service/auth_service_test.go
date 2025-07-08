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

func TestAuthService_CreateTokens(t *testing.T) {
	repo := new(mockRepo)
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "")

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
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "")

	t.Run("cant delete session", func(t *testing.T) {
		repo.On("DeleteSessionByUserID", "u").Return(errors.New("fail")).Once()
		err := svc.Logout("access", "u")
		assert.ErrorIs(t, err, apperrors.ErrCantDeleteSession)
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo.On("DeleteSessionByUserID", "u").Return(nil).Once()
		err := svc.Logout("access", "u")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestAuthService_RefreshTokens(t *testing.T) {
	repo := new(mockRepo)
	svc := NewAuthService(repo, time.Minute, []byte("secret"), "")
	access, _ := makeJWT("u", time.Minute, []byte("secret"))
	hash, _ := bcrypt.GenerateFromPassword([]byte("refresh"), bcrypt.DefaultCost)
	sess := &sessions.Sessions{UserID: "u", RefreshTokenHash: hash, UserAgent: "ua", IPAddr: "ip"}

	t.Run("invalid access token", func(t *testing.T) {
		_, _, err := svc.RefreshTokens("bad", "refresh", "ua", "ip")
		assert.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		repo.On("GetSessionByUserID", "u").Return((*sessions.Sessions)(nil), apperrors.ErrUserNotFound).Once()
		_, _, err := svc.RefreshTokens(access, "refresh", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("cant get session", func(t *testing.T) {
		repo.On("GetSessionByUserID", "u").Return((*sessions.Sessions)(nil), errors.New("fail")).Once()
		_, _, err := svc.RefreshTokens(access, "refresh", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrCantGetSession)
		repo.AssertExpectations(t)
	})

	t.Run("tokens dont match", func(t *testing.T) {
		repo.On("GetSessionByUserID", "u").Return(sess, nil).Once()
		_, _, err := svc.RefreshTokens(access, "wrong", "ua", "ip")
		assert.ErrorIs(t, err, apperrors.ErrTokensDontMatch)
		repo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		repo.On("GetSessionByUserID", "u").Return(sess, nil).Once()
		repo.On("UpdateRefreshTokenByUserID", "u", mock.Anything).Return(nil).Once()
		newAccess, newRefresh, err := svc.RefreshTokens(access, "refresh", "ua", "ip")
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccess)
		assert.NotEmpty(t, newRefresh)
		repo.AssertExpectations(t)
	})
}
