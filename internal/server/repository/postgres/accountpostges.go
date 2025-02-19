package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type accountPostgres struct {
	db *sqlx.DB
}

func NewAccountPostgres(db *sqlx.DB) *accountPostgres {
	return &accountPostgres{
		db: db,
	}
}

func (ap *accountPostgres) SaveAccount(ctx context.Context) (int, error) {
	return 0, nil
}

func (ap *accountPostgres) LoadAccountID(ctx context.Context, login, passwordHash string) (int, error) {
	row := ap.db.QueryRowContext(ctx, `

	SELECT account_id FROM account
	WHERE login=$1 AND password_hash=$2;

	`, login, passwordHash)

	var result int

	err := row.Scan(&result)
	if err != nil {
		return -1, err
	}

	return result, nil
}
