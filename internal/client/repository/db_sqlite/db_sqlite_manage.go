package db_sqlite

import (
	"github.com/golovanevvs/confidant/internal/client/model"
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
