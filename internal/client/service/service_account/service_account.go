package service_account

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportAccount interface {
	CreateAccount(ctx context.Context, email, password string) (trResponse *model.TrResponse, err error)
	Login(ctx context.Context, email, password string) (trResponse *model.TrResponse, err error)
}

type IRepositoryAccount interface {
	SaveAccount(ctx context.Context, accountID int, email string, passwordHash []byte) (err error)
	LoadAccountID(ctx context.Context, email string, passwordHash []byte) (accountID int, err error)
	LoadEmail(ctx context.Context, accountID int) (email string, err error)
	SaveActiveAccount(ctx context.Context, accountID int, refreshTokenString string) (err error)
	LoadActiveAccount(ctx context.Context) (accountID int, refreshTokenstring string, err error)
	DeleteActiveAccount(ctx context.Context) (err error)
}

type ServiceAccount struct {
	tr ITransportAccount
	rp IRepositoryAccount
}

func New(tr ITransportAccount, rp IRepositoryAccount) *ServiceAccount {
	return &ServiceAccount{
		tr: tr,
		rp: rp,
	}
}
