package repo

import (
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}
