package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type myPostgres struct {
	db *sqlx.DB
}

func NewMyPostgres(db *sqlx.DB) *myPostgres {
	return &myPostgres{
		db: db,
	}
}

func (p *myPostgres) DoIt(ctx context.Context) (int, error) {
	return 0, nil
}
