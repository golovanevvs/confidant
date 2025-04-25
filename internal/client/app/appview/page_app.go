package appview

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

type pageApp struct {
	app        *tview.Application
	colorTitle tcell.Color
}

func newPageApp() *pageApp {
	return &pageApp{
		app:        tview.NewApplication(),
		colorTitle: tcell.ColorYellow,
	}
}

func (av *appView) vApp() (err error) {
	action := "run"

	//! set root
	av.v.pageApp.app.SetRoot(av.v.pageMain.mainGrid, true)

	//! get active account
	var loginAtStartErr error
	av.account.ID, av.account.Email, av.refreshToken, loginAtStartErr = av.sv.LoginAtStart(context.Background())

	//! switch group page or login page
	if loginAtStartErr == nil {
		// clearing messages
		av.vClearMessages()

		// setting active account to status bar
		av.v.pageMain.statusBar.cellActiveAccount.SetText(fmt.Sprintf("[green]%s", av.account.Email))

		// switching to groups page
		av.aPageGroupsSwitch()
	} else {
		// clearing messages
		av.vClearMessages()

		// switching to login page
		av.vPageLoginSwitch()
	}

	//! run
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
