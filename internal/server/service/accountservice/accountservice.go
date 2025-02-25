package accountservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/golovanevvs/confidant/internal/server/repository"
)

type accountService struct {
	rp repository.IAccountRepository
}

func NewAccountService(accountRp repository.IAccountRepository) *accountService {
	return &accountService{
		rp: accountRp,
	}
}

func (sv *accountService) CreateAccount(ctx context.Context, account model.Account) (int, error) {
	// e-mail validation
	if !sv.validateEmail(account.Email) {
		return -1, errors.New("e-mail validation error")
	}

	// password hashing
	account.PasswordHash = sv.genPasswordHash(account.Password)

	// DB: saving a new account
	accountID, err := sv.rp.SaveAccount(ctx, account)
	if err != nil {
		action := "create account"
		return -1, fmt.Errorf("%s: %s: %w", customerrors.AccountServiceErr, action, err)
	}

	return accountID, nil
}
