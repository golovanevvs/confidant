package db_sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
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

func (rp *sqliteAccount) SaveAccount(ctx context.Context, accountID int, email string, passwordHash []byte, refreshToken string) (err error) {
	action := "save account"

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO account
			(id, email, password_hash, token)
		VALUES
			(?, ?, ?, ?)

	`, accountID, email, passwordHash, refreshToken)
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

func (rp *sqliteAccount) LoadAccountID(ctx context.Context, email string, passwordHash []byte) (accountID int, err error) {
	action := "load account ID"

	row := rp.db.QueryRowContext(ctx, `

		SELECT id FROM account
		WHERE email=?;

	`, email)

	if err = row.Scan(&accountID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return -1, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBEmailNotFound401,
				err,
			)
		default:
			return -1, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
	}

	row = rp.db.QueryRowContext(ctx, `

		SELECT password_hash FROM account
		WHERE id = ?;

	`, accountID)

	var dbPasswordHash []byte

	if err = row.Scan(&dbPasswordHash); err != nil {
		return -1, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	if !bytes.Equal(dbPasswordHash, passwordHash) {
		return -1, fmt.Errorf(
			"%s: %s: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBWrongPassword401,
		)
	}

	return accountID, nil
}

func (rp *sqliteAccount) LoadEmail(ctx context.Context, accountID int) (email string, err error) {
	action := "load email"

	row := rp.db.QueryRowContext(ctx, `

		SELECT email FROM account
		WHERE id=?;

	`, accountID)

	if err = row.Scan(&email); err != nil {
		return "", fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return email, nil
}

func (rp *sqliteAccount) SaveActiveAccount(ctx context.Context, accountID int) (err error) {
	action := "save active account"

	_, err = rp.db.ExecContext(ctx, `

		INSERT OR REPLACE INTO active_account
			(account_id)
		VALUES
			(?)

	`, accountID)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrSaveActiveAccount,
			err,
		)
	}

	return
}

func (rp *sqliteAccount) LoadActiveAccount(ctx context.Context) (accountID int, refreshTokenstring string, err error) {
	action := "load active account"

	row := rp.db.QueryRowContext(ctx, `

		SELECT
			account_id, token
		FROM
			active_account
		INNER JOIN
			account
		ON
			active_account.account_id = account.id

	`)
	err = row.Scan(&accountID, &refreshTokenstring)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, "", fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrNoActiveAccount,
				err,
			)
		} else {
			return -1, "", fmt.Errorf(
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

func (rp *sqliteAccount) DeleteActiveAccount(ctx context.Context) (err error) {
	action := "delete active account"

	_, err = rp.db.ExecContext(ctx, `

		DELETE FROM active_account

	`)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDeleteActiveAccount,
			err,
		)
	}

	return
}
