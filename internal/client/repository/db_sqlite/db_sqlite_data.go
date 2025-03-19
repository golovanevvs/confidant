package db_sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type sqliteData struct {
	db *sqlx.DB
}

func NewSQLiteData(db *sqlx.DB) *sqliteData {
	return &sqliteData{
		db: db,
	}
}

func (rp *sqliteData) AddGroup(ctx context.Context, accountID int, email string, title string) (err error) {
	action := "add group"

	_, err = rp.db.ExecContext(ctx, `

		INSERT INTO groups
			(title, account_id)
		VALUES
			(?, ?)

	`, title, accountID)

	return nil
}
