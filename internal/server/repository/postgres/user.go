package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type userPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *userPostgres {
	return &userPostgres{
		db: db,
	}
}

func (p *userPostgres) SaveUser(ctx context.Context) (int, error) {
	return 0, nil
}
