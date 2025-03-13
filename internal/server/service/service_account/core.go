package service_account

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (sv *ServiceAccount) CreateAccount(ctx context.Context, account model.Account) (accountID int, err error) {
	action := "create account"

	// password hashing
	account.PasswordHash, err = sv.genHash(account.Password)
	if err != nil {
		return -1, fmt.Errorf("%s: %s: %w", customerrors.AccountServiceErr, action, err)
	}

	// DB: saving a new account
	accountID, err = sv.rp.SaveAccount(ctx, account)
	if err != nil {
		return -1, fmt.Errorf("%s: %s: %w", customerrors.AccountServiceErr, action, err)
	}

	return accountID, nil
}
