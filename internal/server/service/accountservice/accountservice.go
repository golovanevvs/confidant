package accountservice

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IAccountRepository interface {
	SaveAccount(ctx context.Context, account model.Account) (int, error)
	LoadAccountID(ctx context.Context, email, passwordHash string) (int, error)
}

type AccountService struct {
	Rp IAccountRepository
}

func New(accountRp IAccountRepository) *AccountService {
	return &AccountService{
		Rp: accountRp,
	}
}
