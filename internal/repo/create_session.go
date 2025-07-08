package repo

import (
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/Turalchik/authentication-service/internal/entities/sessions"
)

func (repo *Repo) CreateSession(session *sessions.Sessions) error {
	sb := psql.Insert("sessions").
		Columns("user_id", "refresh_token_hash", "user_agent", "ip_addr").
		Values(session.UserID, session.RefreshTokenHash, session.UserAgent, session.IPAddr)

	query, args, err := sb.ToSql()
	if err != nil {
		return apperrors.ErrCantBuildSQLQuery
	}

	if _, err = repo.db.Exec(query, args...); err != nil {
		return apperrors.ErrCantExecSQLQuery
	}
	return nil
}
