package postgres

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type postgresManage struct {
	db *sqlx.DB
}

func NewPostgresManage(db *sqlx.DB) *postgresManage {
	return &postgresManage{
		db: db,
	}
}

func (rp *postgresManage) CloseDB() error {
	err := rp.db.Close()
	if err != nil {
		return fmt.Errorf("error when closing the DB: %s", err.Error())
	}
	return nil
}

func (rp *postgresManage) Ping() (err error) {
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
