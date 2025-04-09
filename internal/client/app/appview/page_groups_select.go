package appview

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageGroupsSelect struct {
	gridSelectButtons *tview.Grid
	buttonSync        *tview.Button
	buttonNew         *tview.Button
	buttonSettings    *tview.Button
	buttonDelete      *tview.Button
	buttonLogout      *tview.Button
	buttonExit        *tview.Button
	grid              *tview.Grid
	inputCapture      func(event *tcell.EventKey) *tcell.EventKey
	page              *tview.Pages
}

func newPageGroupsSelect() *pageGroupsSelect {
	return &pageGroupsSelect{
		gridSelectButtons: tview.NewGrid(),
		buttonSync:        tview.NewButton("Синхронизировать"),
		buttonNew:         tview.NewButton("Создать группу"),
		buttonSettings:    tview.NewButton("Настроить группу"),
		buttonDelete:      tview.NewButton("Удалить группу"),
		buttonLogout:      tview.NewButton("Выйти из аккаунта"),
		buttonExit:        tview.NewButton("Выход"),
		grid:              tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vGroupsSelect() {
	//! "Синхронизировать"
	av.v.pageGroups.pageGroupsSelect.buttonSync.SetSelectedFunc(func() {
		syncResp, err := av.sv.SyncGroups(context.Background(), av.accessToken, av.account.Email)

		// error
		if err != nil {
			av.v.pageMain.statusBar.cellResponseStatus.SetText("")
			av.v.pageMain.messageBoxL.SetText("[red]Ошибка.")
			av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))

			// no error
		} else {
			// setting status
			if syncResp.HTTPStatusCode == 200 {
				av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[green]%s", syncResp.HTTPStatus))
			} else {
				av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[yellow]%s", syncResp.HTTPStatus))
				av.v.pageMain.messageBoxL.SetText("[red]Ошибка.")
				av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", syncResp.Error))
				av.aPageGroupsSwitch()
			}
		}
	})

	//! "Создать группу"
	av.v.pageGroups.pageGroupsSelect.buttonNew.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("add_group_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsAddGroup.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.inputName)
		av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.inputName.SetText("")
	})

	//! "Настроить группу"
	av.v.pageGroups.pageGroupsSelect.buttonSettings.SetSelectedFunc(func() {
		av.aPageGroupsEditEmailsSwitch()
	})

	//! "Удалить группу"
	// av.v.pageGroups.PageSelect.ButtonDelete = tview.NewButton("Удалить группу")

	//! "Выйти из аккаунта"
	av.v.pageGroups.pageGroupsSelect.buttonLogout.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.pages.SwitchToPage("login_page")
		// focus
		av.v.pageApp.app.SetInputCapture(av.v.pageLogin.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputEmail)
		// status bar
		av.v.pageMain.statusBar.cellActiveAccount.SetText("")
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		// messageBox
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
		// form
		av.v.pageLogin.form.inputEmail.SetText("")
		av.v.pageLogin.form.inputPassword.SetText("")

		// deleting active account
		err := av.sv.Logout(context.Background())
		if err != nil {
			av.v.pageMain.messageBoxL.SetText("[red]Критическая ошибка!")
			av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
		}
		av.accessToken = ""
		av.refreshToken = ""
		av.account.ID = -1
		av.account.Email = ""
	})

	//! "Выход"
	av.v.pageGroups.pageGroupsSelect.buttonExit.SetSelectedFunc(func() {
		av.v.pageApp.app.Stop()
	})

	//! MainButtonsGrid
	av.v.pageGroups.pageGroupsSelect.grid.
		SetRows(1, 1, 1, 1, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageGroupsSelect.buttonSync, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsSelect.buttonNew, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsSelect.buttonSettings, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsSelect.buttonDelete, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsSelect.buttonLogout, 5, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsSelect.buttonExit, 6, 0, 1, 1, 0, 0, true)

	//! InputCapture select page
	av.v.pageGroups.pageGroupsSelect.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.listGroups:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsSelect.buttonSync)
			case av.v.pageGroups.pageGroupsSelect.buttonSync:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsSelect.buttonNew)
			case av.v.pageGroups.pageGroupsSelect.buttonNew:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsSelect.buttonSettings)
			case av.v.pageGroups.pageGroupsSelect.buttonSettings:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsSelect.buttonDelete)
			case av.v.pageGroups.pageGroupsSelect.buttonDelete:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsSelect.buttonLogout)
			case av.v.pageGroups.pageGroupsSelect.buttonLogout:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsSelect.buttonExit)
			case av.v.pageGroups.pageGroupsSelect.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
			}
			return nil
		}
		return event
	}
}

func (av *appView) aPageGroupsEditEmailsSwitch() {
	av.vClearMessages()
	av.v.pageGroups.pages.SwitchToPage("edit_emails_page")
	av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsEditEmails.inputCapture)
	av.v.pageApp.app.SetFocus(av.v.pageGroups.listEmails)
	av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.inputEmail.SetText("")
}
