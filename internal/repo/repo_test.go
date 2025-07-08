package repo

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/Turalchik/authentication-service/internal/entities/sessions"
	"github.com/jmoiron/sqlx"
)

func setupDataBase(t *testing.T) (*Repo, sqlmock.Sqlmock, func(), error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, err
	}

	db := sqlx.NewDb(sqlDB, "sqlmock")
	repoObj := NewRepo(db)

	return repoObj, mock, func() { db.Close() }, nil
}

func TestRepo_GetSessionByUserID(t *testing.T) {
	repo, mock, closer, err := setupDataBase(t)
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer closer()

	expectQuery := regexp.QuoteMeta("SELECT * FROM sessions WHERE user_id = $1")
	expectSession := sessions.Sessions{
		UserID:           "user_id_test",
		RefreshTokenHash: []byte("refresh_token_hash_test"),
		UserAgent:        "user_agent_test",
		IPAddr:           "ip_addr_test",
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"user_id", "refresh_token_hash", "user_agent", "ip_addr"}).
			AddRow(expectSession.UserID, expectSession.RefreshTokenHash, expectSession.UserAgent, expectSession.IPAddr)

		mock.
			ExpectQuery(expectQuery).
			WithArgs("user_id_test").
			WillReturnRows(rows)

		session, err := repo.GetSessionByUserID("user_id_test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if session == nil {
			t.Fatal("expected non-nil session")
		}
		if session.UserID != expectSession.UserID {
			t.Errorf("UserID = %s; want %s", session.UserID, expectSession.UserID)
		}
		if string(session.RefreshTokenHash) != string(expectSession.RefreshTokenHash) {
			t.Errorf("RefreshTokenHash = %s; want %s", string(session.RefreshTokenHash), string(expectSession.RefreshTokenHash))
		}
		if session.UserAgent != expectSession.UserAgent {
			t.Errorf("UserAgent = %s; want %s", session.UserAgent, expectSession.UserAgent)
		}
		if session.IPAddr != expectSession.IPAddr {
			t.Errorf("IPAddr = %s; want %s", session.IPAddr, expectSession.IPAddr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"user_id", "refresh_token_hash", "user_agent", "ip_addr"})

		mock.
			ExpectQuery(expectQuery).
			WithArgs("user_id_test").
			WillReturnRows(rows)

		session, err := repo.GetSessionByUserID("user_id_test")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if session != nil {
			t.Fatal("expected nil session")
		}
		if !errors.Is(err, apperrors.ErrUserNotFound) {
			t.Fatalf("expected sql no rows error, got: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestRepo_CreateSession(t *testing.T) {
	repo, mock, closer, err := setupDataBase(t)
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer closer()

	expectQuery := regexp.QuoteMeta("INSERT INTO sessions (user_id,refresh_token_hash,user_agent,ip_addr) VALUES ($1,$2,$3,$4)")
	sess := &sessions.Sessions{
		UserID:           "user_id_test",
		RefreshTokenHash: []byte("refresh_token_hash_test"),
		UserAgent:        "user_agent_test",
		IPAddr:           "ip_addr_test",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(expectQuery).
			WithArgs(sess.UserID, sess.RefreshTokenHash, sess.UserAgent, sess.IPAddr).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateSession(sess)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("sql error", func(t *testing.T) {
		mock.ExpectExec(expectQuery).
			WithArgs(sess.UserID, sess.RefreshTokenHash, sess.UserAgent, sess.IPAddr).
			WillReturnError(errors.New("db error"))

		err := repo.CreateSession(sess)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestRepo_DeleteSessionByUserID(t *testing.T) {
	repo, mock, closer, err := setupDataBase(t)
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer closer()

	expectQuery := regexp.QuoteMeta("DELETE FROM sessions WHERE user_id = $1")

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(expectQuery).
			WithArgs("user_id_test").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.DeleteSessionByUserID("user_id_test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("sql error", func(t *testing.T) {
		mock.ExpectExec(expectQuery).
			WithArgs("user_id_test").
			WillReturnError(errors.New("db error"))
		err := repo.DeleteSessionByUserID("user_id_test")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestRepo_UpdateRefreshTokenByUserID(t *testing.T) {
	repo, mock, closer, err := setupDataBase(t)
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer closer()

	expectQuery := regexp.QuoteMeta("UPDATE sessions SET refresh_token_hash = $1 WHERE user_id = $2")

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(expectQuery).
			WithArgs("refresh_token_hash_test", "user_id_test").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.UpdateRefreshTokenByUserID("user_id_test", "refresh_token_hash_test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("sql error", func(t *testing.T) {
		mock.ExpectExec(expectQuery).
			WithArgs("refresh_token_hash_test", "user_id_test").
			WillReturnError(errors.New("db error"))
		err := repo.UpdateRefreshTokenByUserID("user_id_test", "refresh_token_hash_test")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
