package appview

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
	"go.uber.org/zap"
)

type IServiceAccount interface {
	CreateAccount(ctx context.Context, email, password string) (registerAccountResp *model.RegisterAccountResp, err error)
	GetAccessToken(ctx context.Context, refreshTokenString string) (accessTokenString string, err error)
	// Login(email, password string) (accountID int, err error)
	// ChangePassword(email, password, newPassword string) error
	//GetUser(login string) (string, error)
}

type IServiceGroups interface {
	GetGroups(ctx context.Context)
}

type IServiceStatusServer interface {
	GetServerStatus(ctx context.Context) (statusResp *model.StatusResp, err error)
}

type IService interface {
	IServiceAccount
	IServiceStatusServer
}

type view struct {
	pageApp      *pageApp
	pageMain     *pageMain
	pageLogin    *pageLogin
	pageRegister *pageRegister
	pageGroups   *pageGroups
}

type appView struct {
	sv          IService
	lg          *zap.SugaredLogger
	v           view
	accessToken string
}

func New(sv IService, lg *zap.SugaredLogger) *appView {
	return &appView{
		v: view{
			pageApp:      newPageApp(),
			pageMain:     newPageMain(),
			pageLogin:    newPageLogin(),
			pageRegister: newPageRegister(),
			pageGroups:   newPageGroups(),
		},
		sv: sv,
		lg: lg,
	}
}

func (av *appView) Run() error {

	av.vMain()
	av.vLogin()
	av.vRegister()
	av.vGroups()

	return av.vApp()
}
