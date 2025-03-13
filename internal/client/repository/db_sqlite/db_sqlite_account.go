package db_sqlite

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
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

func (rp *sqliteAccount) SaveAccount(email string, passwordHash []byte, refreshTokenString string) error {
	action := "save refresh token"

	ctx := context.Background()

	_, err := rp.db.ExecContext(ctx, `
	
		INSERT INTO account
			(email, password_hash, refresh_token)
		VALUES
			($1, $2, $3)

	`, email, passwordHash, refreshTokenString)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrSaveRefreshToken,
			err,
		)
	}

	return nil
}

func (rp *sqliteAccount) LoadAccountID(email, passwordHash string) (int, error) {
	return 0, nil
}

// func (rp *sqliteAccount) SaveRefreshToken(email, refreshTokenString string) error {
// 	action := "save refresh token"

// 	ctx := context.Background()

// 	_, err := rp.db.ExecContext(ctx, `

// 		INSERT INTO account
// 			(email, refresh_token)
// 		VALUES
// 			($1, $2)

// 	`, email, refreshTokenString)
// 	if err != nil {
// 		return fmt.Errorf(
// 			"%s: %s: %w: %w",
// 			customerrors.DBErr,
// 			action,
// 			customerrors.ErrSaveRefreshToken,
// 			err,
// 		)
// 	}

// 	return nil
// }
