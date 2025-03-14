package db_sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (rp *sqliteAccount) SaveAccount(ctx context.Context, email string, passwordHash []byte) (err error) {
	action := "save account"

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO account
			(email, password_hash)
		VALUES
			($1, $2, $3)

	`, email, passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBBusyEmail409,
				err,
			)
		} else {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
	}

	return nil
}

func (rp *sqliteAccount) LoadAccountID(ctx context.Context, email, passwordHash string) (int, error) {
	return 0, nil
}

func (rp *sqliteAccount) LoadActiveRefreshToken(ctx context.Context) (refreshTokenstring string, err error) {
	action := "load active refresh token"

	row := rp.db.QueryRowContext(ctx, `

		SELECT token FROM refresh_token

	`)
	err = row.Scan(&refreshTokenstring)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrNoActiveRefreshToken,
				err,
			)
		} else {
			return "", fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
	}

	return
}

func (rp *sqliteAccount) SaveRefreshToken(ctx context.Context, refreshTokenString string) (err error) {
	action := "save active refresh token"

	_, err = rp.db.ExecContext(ctx, `

		INSERT INTO refresh_token
			(token)
		VALUES
			($1)

	`, refreshTokenString)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrSaveActiveRefreshToken,
			err,
		)
	}

	return
}

func (rp *sqliteAccount) DeleteActiveRefreshToken(ctx context.Context) (err error) {
	action := "delete active refresh token"

	_, err = rp.db.ExecContext(ctx, `

		DELETE FROM active_refresh_token

	`)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDeleteActiveRefreshToken,
			err,
		)
	}

	return
}
