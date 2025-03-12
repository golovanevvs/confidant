package sqlite

import (
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/jmoiron/sqlx"
)

type sqliteAccount struct {
	db *sqlx.DB
}

func NewSQLiteAccount(db *sqlx.DB) *sqliteAccount {
	return &sqliteAccount{
		db: db,
	}
}

func (rp *sqliteAccount) SaveAccount(account model.Account) (int, error) {
	return 0, nil
}

func (rp *sqliteAccount) LoadAccountID(email, passwordHash string) (int, error) {
	return 0, nil
}
