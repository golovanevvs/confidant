package db_sqlite

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
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

func (rp *sqliteAccount) GetServerStatus() (statusResp *model.StatusResp, err error) {
	return nil, nil
}

func (rp *sqliteAccount) CloseDB() error {
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
