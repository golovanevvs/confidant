package appview

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
	"go.uber.org/zap"
)

type IServiceAccount interface {
	CreateAccount(ctx context.Context, email, password string) (registerAccountResp *model.AccountResp, err error)
	GetAccessToken(ctx context.Context, refreshTokenString string) (accessTokenString string, err error)
	Login(ctx context.Context, email, password string) (registerAccountResp *model.AccountResp, err error)
	LoginAtStart(ctx context.Context) (accountID int, email string, refreshTokenString string, err error)
	Logout(ctx context.Context) (err error)
}

type IServiceGroups interface {
	// GetGroups(ctx context.Context)
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
}

type IServiceManage interface {
	GetServerStatus(ctx context.Context) (statusResp *model.StatusResp, err error)
}

type IService interface {
	IServiceAccount
	IServiceManage
	IServiceGroups
}

type view struct {
	pageApp      *pageApp
	pageMain     *pageMain
	pageLogin    *pageLogin
	pageRegister *pageRegister
	pageGroups   *pageGroups
	pageData     *pageData
}

type appView struct {
	sv           IService
	lg           *zap.SugaredLogger
	v            view
	accessToken  string
	refreshToken string
	account      model.Account
}

func New(sv IService, lg *zap.SugaredLogger) *appView {
	return &appView{
		v: view{
			pageApp:      newPageApp(),
			pageMain:     newPageMain(),
			pageLogin:    newPageLogin(),
			pageRegister: newPageRegister(),
			pageGroups:   newPageGroups(),
			pageData:     newPageData(),
		},
		sv:           sv,
		lg:           lg,
		accessToken:  "",
		refreshToken: "",
		account: model.Account{
			ID:    -1,
			Email: "",
		},
	}
}

func (av *appView) Run() error {

	av.vMain()
	av.vLogin()
	av.vRegister()
	av.vGroups()
	av.vGroupsSelect()
	av.vGroupsAddGroup()
	av.vGroupsEditEmails()
	av.vData()
	av.vDataSelectType()
	av.vDataViewNote()
	av.vDataViewPass()
	av.VDataViewCard()
	av.VDataViewFile()
	av.vDataAddNote()
	av.vDataAddPass()
	av.vDataAddCard()
	av.VDataAddFile()

	return av.vApp()
}
