package database

import (
	"fmt"
	"github.com/Turalchik/authentication-service/internal/apperrors"
	"github.com/jmoiron/sqlx"
	"os"
)

func NewDatabase(dsn string, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, apperrors.ErrCantOpenDatabase
	}
	return db, nil
}

func NewPostgresDSN() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbUser, dbPassword, dbHost, dbPort, dbName)
}
