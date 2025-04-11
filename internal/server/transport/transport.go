package transport

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IServiceAccount interface {
	CreateAccount(ctx context.Context, account model.Account) (accountID int, err error)
	Login(ctx context.Context, account model.Account) (accountID int, err error)
	BuildAccessJWTString(ctx context.Context, accountID int) (accessTokenString string, err error)
	BuildRefreshJWTString(ctx context.Context, accountID int) (refreshTokenString string, err error)
	GetAccountIDFromJWT(tokenString string) (int, error)
	RefreshAccessJWT(ctx context.Context, refreshToken string) (accessTokenString string, err error)
}

type IServiceManage interface {
	PingDB() (err error)
}

type IServiceGroups interface {
	GetGroupIDs(ctx context.Context, accountID int) (groupIDs []int, err error)
	GetGroups(ctx context.Context, accountID int, groupIDs []int) (groups []model.Group, err error)
	AddGroups(ctx context.Context, groups []model.Group) (groupIDs map[int]int, err error)
}

type IServiceData interface {
	GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error)
}

type IService interface {
	IServiceAccount
	IServiceManage
	IServiceGroups
}
