package transport

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, account model.Account) (accountID int, err error)
	BuildAccessJWTString(ctx context.Context, accountID int) (accessTokenString string, err error)
	BuildRefreshJWTString(ctx context.Context, accountID int) (refreshTokenString string, err error)
}

type IOtherService interface {
	DoSomething() error
}

type IService interface {
	IAccountService
	IOtherService
}
