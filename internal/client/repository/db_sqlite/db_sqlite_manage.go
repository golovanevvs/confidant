package db_sqlite

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type sqliteManage struct {
	db *sqlx.DB
}

func NewSQLiteManage(db *sqlx.DB) *sqliteManage {
	return &sqliteManage{
		db: db,
	}
}

func (rp *sqliteAccount) CloseDB(ctx context.Context) error {
	action := "close DB"

	if err := rp.db.Close(); err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.DBErr,
			action,
			err,
		)
	}

	return nil
}
