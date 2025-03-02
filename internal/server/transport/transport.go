package transport

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account model.Account) (accountID int, err error)
	BuildJWTString(ctx context.Context, accountID int) (tokenString string, err error)
}

type IOtherService interface {
	DoSomething() error
}

type IService interface {
	IAccountService
	IOtherService
}
