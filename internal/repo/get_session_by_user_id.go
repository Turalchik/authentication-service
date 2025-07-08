package repo

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/Turalchik/authentication-service/internal/entities/sessions"
)

func (repo *Repo) GetSessionByUserID(userID string) (*sessions.Sessions, error) {
	sb := psql.Select("*").
		From("sessions").
		Where(sq.Eq{"user_id": userID})

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, apperrors.ErrCantBuildSQLQuery
	}

	session := &sessions.Sessions{}
	if err = repo.db.Get(session, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, apperrors.ErrCantExecSQLQuery
	}

	return session, nil
}
