package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/jmoiron/sqlx"
)

type AccountPostgres struct {
	db *sqlx.DB
}

func NewAccountPostgres(db *sqlx.DB) *AccountPostgres {
	return &AccountPostgres{
		db: db,
	}
}

func (rp *AccountPostgres) SaveAccount(ctx context.Context, account model.Account) (int, error) {
	action := "save account"

	row := rp.db.QueryRowContext(ctx, `

		INSERT INTO account
			(email, password_hash)
		VALUES
			($1, $2)
		RETURNING id;

	`, account.Email, account.PasswordHash)

	var accountID int
	if err := row.Scan(&accountID); err != nil {
		switch {
		case strings.Contains(err.Error(), " 23505"):
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBBusyEmail409, err)
		default:
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBInternalError500, err)
		}
	}

	return accountID, nil
}

func (rp *AccountPostgres) LoadAccountID(ctx context.Context, email, passwordHash string) (int, error) {
	action := "load account ID"

	row := rp.db.QueryRowContext(ctx, `

	SELECT id FROM account
	WHERE email=$1;

	`, email)

	var accountID int

	if err := row.Scan(&accountID); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBEmailNotFound401, err)
		default:
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBInternalError500, err)
		}
	}

	row = rp.db.QueryRowContext(ctx, `

	SELECT password_hash FROM account
	WHERE id = $1;

	`, accountID)

	var dbPasswordHash string

	if err := row.Scan(&dbPasswordHash); err != nil {
		return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBInternalError500, err)
	}

	if dbPasswordHash != passwordHash {
		return -1, fmt.Errorf("%s: %s: %w", customerrors.DBErr, action, customerrors.ErrDBWrongPassword401)
	}

	return accountID, nil
}
