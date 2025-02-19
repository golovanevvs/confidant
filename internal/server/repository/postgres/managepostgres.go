package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type managePostgres struct {
	db *sqlx.DB
}

func NewManagePostgres(db *sqlx.DB) *managePostgres {
	return &managePostgres{
		db: db,
	}
}

func (mp *managePostgres) CloseDB() error {
	err := mp.db.Close()
	if err != nil {
		return fmt.Errorf("error when closing the DB: %s", err.Error())
	}
	return nil
}
