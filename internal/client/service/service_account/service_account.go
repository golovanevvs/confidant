package service_account

import (
	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransportAccount interface {
	CreateAccount(email, password string) (trResponse *model.TrResponse, err error)
}

type IRepositoryAccount interface {
	SaveAccount(email string, passwordHash []byte, refreshTokenString string) error
	LoadAccountID(email, passwordHash string) (int, error)
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
