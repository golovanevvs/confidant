package appview

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

type pageApp struct {
	app *tview.Application
}

func newPageApp() *pageApp {
	return &pageApp{
		app: tview.NewApplication(),
	}
}

func (av *appView) vApp() error {
	action := "run"

	var loginAtStartErr error

	av.v.pageApp.app.SetRoot(av.v.pageMain.mainGrid, true)

	av.account.ID, av.account.Email, av.refreshToken, loginAtStartErr = av.sv.LoginAtStart(context.Background())
	if loginAtStartErr == nil {
		av.v.pageMain.statusBar.cellActiveAccount.SetText(fmt.Sprintf("[green]%s", av.account.Email))
		av.v.pageMain.pages.SwitchToPage("groups_page")
		av.v.pageGroups.pagesSelEd.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
	} else {
		av.v.pageMain.pages.SwitchToPage("login_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageLogin.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputEmail)
	}

	if err := av.v.pageApp.app.Run(); err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientAppViewErr,
			action,
			customerrors.ErrRunAppView,
			err)
	}
	return nil
}
