package service_account

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IRepositoryAccount interface {
	SaveAccount(ctx context.Context, account model.Account) (int, error)
	LoadAccountID(ctx context.Context, email string, passwordHash []byte) (int, error)
	SaveRefreshTokenHash(ctx context.Context, accountID int, refreshTokenHash []byte) error
}

type ServiceAccount struct {
	rp IRepositoryAccount
}

func New(accountRp IRepositoryAccount) *ServiceAccount {
	return &ServiceAccount{
		rp: accountRp,
	}
}
