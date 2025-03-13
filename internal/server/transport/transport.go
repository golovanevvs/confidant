package transport

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IServiceAccount interface {
	CreateAccount(ctx context.Context, account model.Account) (accountID int, err error)
	BuildAccessJWTString(ctx context.Context, accountID int) (accessTokenString string, err error)
	BuildRefreshJWTString(ctx context.Context, accountID int) (refreshTokenString string, err error)
	GetAccountIDFromJWT(tokenString string) (int, error)
}

type IServiceManage interface {
	PingDB() (err error)
}

type IService interface {
	IServiceAccount
	IServiceManage
}
