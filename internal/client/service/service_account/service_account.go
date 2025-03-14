package service_account

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportAccount interface {
	CreateAccount(ctx context.Context, email, password string) (trResponse *model.TrResponse, err error)
}

type IRepositoryAccount interface {
	SaveAccount(ctx context.Context, email string, passwordHash []byte) (err error)
	LoadAccountID(ctx context.Context, email, passwordHash string) (int, error)
	LoadActiveRefreshToken(ctx context.Context) (refreshTokenstring string, err error)
	SaveRefreshToken(ctx context.Context, refreshTokenString string) (err error)
	DeleteActiveRefreshToken(ctx context.Context) (err error)
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
