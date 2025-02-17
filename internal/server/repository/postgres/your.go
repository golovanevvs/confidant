package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type yourPostgres struct {
	db *sqlx.DB
}

func NewYourPostgres(db *sqlx.DB) *yourPostgres {
	return &yourPostgres{
		db: db,
	}
}

func (p *yourPostgres) CrashIt(ctx context.Context) (int, error) {
	return 0, nil
}
