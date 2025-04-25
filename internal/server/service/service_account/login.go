package service_account

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (sv *ServiceAccount) Login(ctx context.Context, account model.Account) (accountID int, err error) {
	action := "login"

	// password hashing
	account.PasswordHash, err = sv.genHash(account.Password)
	if err != nil {
		return -1, fmt.Errorf(
			"%s: %s: %w",
			customerrors.AccountServiceErr,
			action,
			err,
		)
	}

	accountID, err = sv.rp.LoadAccountID(ctx, account.Email, account.PasswordHash)
	if err != nil {
		return -1, fmt.Errorf(
			"%s: %s: %w",
			customerrors.AccountServiceErr,
			action,
			err,
		)
	}

	return accountID, nil
}
