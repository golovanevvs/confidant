package service_account

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IRepositoryAccount interface {
	SaveAccount(ctx context.Context, account model.Account) (int, error)
	LoadAccountID(ctx context.Context, email, passwordHash string) (int, error)
}

type ServiceAccount struct {
	rp IRepositoryAccount
}

func New(accountRp IRepositoryAccount) *ServiceAccount {
	return &ServiceAccount{
		rp: accountRp,
	}
}
