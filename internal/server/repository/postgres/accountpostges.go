package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/jmoiron/sqlx"
)

type accountPostgres struct {
	db *sqlx.DB
}

func NewAccountPostgres(db *sqlx.DB) *accountPostgres {
	return &accountPostgres{
		db: db,
	}
}

func (rp *accountPostgres) SaveAccount(ctx context.Context, account model.Account) (int, error) {
	row := rp.db.QueryRowContext(ctx, `
		INSERT INTO account
			(email, password_hash)
		VALUES
			($1, $2)
		RETURNING id;
	`, account.Email, account.PasswordHash)

	var accountID int
	if err := row.Scan(&accountID); err != nil {
		action := "save account"
		switch {
		case strings.Contains(err.Error(), " 23505"):
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBBusyEmail409, err)
		default:
			return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.DBErr, action, customerrors.ErrDBInternalError500, err)
		}
	}

	return accountID, nil
}

func (rp *accountPostgres) LoadAccountID(ctx context.Context, login, passwordHash string) (int, error) {
	row := rp.db.QueryRowContext(ctx, `

	SELECT account_id FROM account
	WHERE login=$1 AND password_hash=$2;

	`, login, passwordHash)

	var result int

	err := row.Scan(&result)
	if err != nil {
		return -1, err
	}

	return result, nil
}
