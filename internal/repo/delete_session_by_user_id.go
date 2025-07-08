package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/Turalchik/authentication-service/internal/apperrors"
)

func (repo *Repo) DeleteSessionByUserID(userID string) error {
	sb := psql.Delete("sessions").
		Where(sq.Eq{"user_id": userID})

	query, args, err := sb.ToSql()
	if err != nil {
		return apperrors.ErrCantBuildSQLQuery
	}

	if _, err = repo.db.Exec(query, args...); err != nil {
		return apperrors.ErrCantExecSQLQuery
	}
	return nil
}
