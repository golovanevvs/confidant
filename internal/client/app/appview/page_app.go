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

func (av *appView) vApp() error {
	action := "run"

	var loginAtStartErr error

	av.v.pageApp.app.SetRoot(av.v.pageMain.mainGrid, true)

	av.account.ID, av.account.Email, av.refreshToken, loginAtStartErr = av.sv.LoginAtStart(context.Background())
	if loginAtStartErr == nil {
		var err error
		av.groups, err = av.sv.GetGroups(context.Background(), av.account.ID)
		if err != nil {
			// error
			av.v.pageMain.statusBar.cellResponseStatus.SetText("")
			av.v.pageMain.messageBoxL.SetText(fmt.Sprintf("[red]Ошибка: %s", err.Error()))
			av.v.pageMain.messageBoxR.Clear()

		} else {
			// updating groups list
			av.v.pageGroups.listGroups.Clear()
			for _, group := range av.groups {
				av.v.pageGroups.listGroups.AddItem(group.Title, "", 0, nil)
			}

			// updating e-mails list
			av.v.pageGroups.listEmails.Clear()
			for _, email := range av.groups[0].Emails {
				av.v.pageGroups.listEmails.AddItem(email, "", 0, nil)
			}

			// set active account to status bar
			av.v.pageMain.statusBar.cellActiveAccount.SetText(fmt.Sprintf("[green]%s", av.account.Email))

			// switch
			av.v.pageMain.pages.SwitchToPage("groups_page")
			av.v.pageGroups.pages.SwitchToPage("select_page")
			av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsSelect.inputCapture)
			av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
		}
	} else {
		// switch
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
