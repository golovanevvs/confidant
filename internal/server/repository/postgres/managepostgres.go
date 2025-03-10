package postgres

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
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

func (rp *ManagePostgres) CloseDB() error {
	err := rp.db.Close()
	if err != nil {
		return fmt.Errorf("error when closing the DB: %s", err.Error())
	}
	return nil
}

func (rp *ManagePostgres) Ping() (err error) {
	action := "ping DB"
	if err := rp.db.Ping(); err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	return nil
}
