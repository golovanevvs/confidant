package accountservice

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (sv *AccountService) CreateAccount(ctx context.Context, account model.Account) (accountID int, err error) {
	action := "create account"

	// password hashing
	account.PasswordHash = sv.genPasswordHash(account.Password)

	// DB: saving a new account
	accountID, err = sv.Rp.SaveAccount(ctx, account)
	if err != nil {
		return -1, fmt.Errorf("%s: %s: %w", customerrors.AccountServiceErr, action, err)
	}

	return accountID, nil
}
