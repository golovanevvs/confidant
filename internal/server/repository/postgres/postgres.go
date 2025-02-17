package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// The postgres constructor.
func New(databaseURI string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	// DB pinging
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("DB ping error: %s", err.Error())
	}

	return db, nil
}
