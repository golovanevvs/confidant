package sqlite

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type SQLite struct {
	*sqliteAccount
	*sqliteManage
}

func New(databasePath string) (*SQLite, error) {
	db, err := sqlx.Open("sqlite3", databasePath)
	if err != nil {
		return nil, err
	}
	return &SQLite{
		sqliteAccount: NewSQLiteAccount(db),
		sqliteManage:  NewSQLiteManage(db),
	}, nil
}
