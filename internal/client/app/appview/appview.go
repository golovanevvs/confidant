package appview

import (
	"github.com/golovanevvs/confidant/internal/client/model"
	"go.uber.org/zap"
)

type IAccountService interface {
	RegisterAccount(email, password string) (registerAccountResp *model.RegisterAccountResp, err error)
	// Login(email, password string) (accountID int, err error)
	// ChangePassword(email, password, newPassword string) error
	//GetUser(login string) (string, error)
}

type IStatusServerService interface {
	GetServerStatus() (statusResp *model.StatusResp, err error)
}

type IService interface {
	IAccountService
	IStatusServerService
}

type view struct {
	pageApp      *pageApp
	pageMain     *pageMain
	pageLogin    *pageLogin
	pageRegister *pageRegister
	pageGroups   *pageGroups
}

type appView struct {
	sv IService
	lg *zap.SugaredLogger
	v  view
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
