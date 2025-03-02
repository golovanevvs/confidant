package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ManagePostgres struct {
	db *sqlx.DB
}

func NewManagePostgres(db *sqlx.DB) *ManagePostgres {

	return &ManagePostgres{
		db: db,
	}
}

func (mp *ManagePostgres) CloseDB() error {
	err := mp.db.Close()
	if err != nil {
		return fmt.Errorf("error when closing the DB: %s", err.Error())
	}
	return nil
}
