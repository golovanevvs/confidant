package appview

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageSelect struct {
	gridSelectButtons *tview.Grid
	buttonSync        *tview.Button
	buttonNew         *tview.Button
	buttonSettings    *tview.Button
	buttonDelete      *tview.Button
	buttonLogout      *tview.Button
	buttonExit        *tview.Button
	grid              *tview.Grid
	inputCapture      func(event *tcell.EventKey) *tcell.EventKey
	Page              *tview.Pages
}

type pageAddGroup struct {
	formAddGroup formAddGroup
	buttonNew    *tview.Button
	buttonExit   *tview.Button
	grid         *tview.Grid
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
	page         *tview.Pages
}

type formAddGroup struct {
	inputName *tview.InputField
	form      *tview.Form
}

type pageEditEmails struct {
	formAddEmail formAddEmail
	buttonAdd    *tview.Button
	buttonDelete *tview.Button
	buttonEхit   *tview.Button
	grid         *tview.Grid
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
	page         *tview.Pages
}

type formAddEmail struct {
	inputEmail *tview.InputField
	form       *tview.Form
}

type pageGroups struct {
	listGroups     *tview.List
	listEmails     *tview.List
	gridMain       *tview.Grid
	pageSelect     pageSelect
	pageAddGroup   pageAddGroup
	pageEditEmails pageEditEmails
	pages          *tview.Pages
}

func newPageGroups() *pageGroups {
	return &pageGroups{
		listGroups: tview.NewList(),
		listEmails: tview.NewList(),
		gridMain:   tview.NewGrid(),
		pageSelect: pageSelect{
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
			Page: tview.NewPages(),
		},
		pageAddGroup: pageAddGroup{
			formAddGroup: formAddGroup{
				inputName: tview.NewInputField(),
				form:      tview.NewForm(),
			},
			buttonNew:  tview.NewButton("Создать группу"),
			buttonExit: tview.NewButton("Назад"),
			grid:       tview.NewGrid(),
			inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
				return event
			},
			page: tview.NewPages(),
		},
		pageEditEmails: pageEditEmails{
			formAddEmail: formAddEmail{
				inputEmail: tview.NewInputField(),
				form:       tview.NewForm(),
			},
			buttonAdd:    tview.NewButton("Добавить e-mail"),
			buttonDelete: tview.NewButton("Удалить e-mail"),
			buttonEхit:   tview.NewButton("Назад"),
			grid:         tview.NewGrid(),
			inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
				return event
			},
			page: tview.NewPages(),
		},
		pages: tview.NewPages(),
	}
}

func (av *appView) vGroups() {
	//! Groups List
	av.v.pageGroups.listGroups.ShowSecondaryText(false)
	av.v.pageGroups.listGroups.SetBorder(true)
	av.v.pageGroups.listGroups.SetHighlightFullLine(true)
	av.v.pageGroups.listGroups.SetTitle(" Список групп ")
	for i := 0; i < 10; i++ {
		av.v.pageGroups.listGroups.AddItem(fmt.Sprintf("Group %d", i), "", 0, nil)
	}

	av.v.pageGroups.listGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	av.v.pageGroups.listGroups.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText)
	})

	//! "select_page"

	//! "Синхронизировать"

	//! "Создать группу"
	av.v.pageGroups.pageSelect.buttonNew.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("add_group_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageAddGroup.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.pageAddGroup.formAddGroup.inputName)
		av.v.pageGroups.pageAddGroup.formAddGroup.inputName.SetText("")
	})

	//! "Настроить группу"
	av.v.pageGroups.pageSelect.buttonSettings.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("edit_emails_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageEditEmails.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listEmails)
		av.v.pageGroups.pageEditEmails.formAddEmail.inputEmail.SetText("")
	})

	//! "Удалить группу"
	// av.v.pageGroups.PageSelect.ButtonDelete = tview.NewButton("Удалить группу")

	//! "Выйти из аккаунта"
	av.v.pageGroups.pageSelect.buttonLogout.SetSelectedFunc(func() {
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
	av.v.pageGroups.pageSelect.buttonExit.SetSelectedFunc(func() {
		av.v.pageApp.app.Stop()
	})

	//! MainButtonsGrid
	av.v.pageGroups.pageSelect.grid.
		SetRows(1, 1, 1, 1, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageSelect.buttonSync, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageSelect.buttonNew, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageSelect.buttonSettings, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageSelect.buttonDelete, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageSelect.buttonLogout, 5, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageSelect.buttonExit, 6, 0, 1, 1, 0, 0, true)

	//! InputCapture select page
	av.v.pageGroups.pageSelect.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.listGroups:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonSync)
			case av.v.pageGroups.pageSelect.buttonSync:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonNew)
			case av.v.pageGroups.pageSelect.buttonNew:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonSettings)
			case av.v.pageGroups.pageSelect.buttonSettings:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonDelete)
			case av.v.pageGroups.pageSelect.buttonDelete:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonLogout)
			case av.v.pageGroups.pageSelect.buttonLogout:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonExit)
			case av.v.pageGroups.pageSelect.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
			}
			return nil
		}
		return event
	}

	//! "add_group_page"

	//! add form
	av.v.pageGroups.pageAddGroup.formAddGroup.form.SetHorizontal(false)
	av.v.pageGroups.pageAddGroup.formAddGroup.form.AddInputField("", "", 0, nil, nil)
	av.v.pageGroups.pageAddGroup.formAddGroup.inputName = av.v.pageGroups.pageAddGroup.formAddGroup.form.GetFormItem(0).(*tview.InputField)

	//! "Создать группу"

	//! "Назад"
	av.v.pageGroups.pageAddGroup.buttonExit.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
	})

	//! "buttons grid"
	av.v.pageGroups.pageAddGroup.grid.
		SetRows(3, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageAddGroup.formAddGroup.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageAddGroup.buttonNew, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageAddGroup.buttonExit, 2, 0, 1, 1, 0, 0, true)

	//! InputCapture add group page
	av.v.pageGroups.pageAddGroup.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.pageAddGroup.formAddGroup.inputName:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageAddGroup.buttonNew)
			case av.v.pageGroups.pageAddGroup.buttonNew:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageAddGroup.buttonExit)
			case av.v.pageGroups.pageAddGroup.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageAddGroup.formAddGroup.inputName)
			}
			return nil
		}
		return event
	}

	//! "edit_emails_page"

	//! add form
	av.v.pageGroups.pageEditEmails.formAddEmail.form.SetHorizontal(false)
	av.v.pageGroups.pageEditEmails.formAddEmail.form.AddInputField("", "", 0, nil, nil)
	av.v.pageGroups.pageEditEmails.formAddEmail.inputEmail = av.v.pageGroups.pageEditEmails.formAddEmail.form.GetFormItem(0).(*tview.InputField)

	//! "Добавить e-mail"
	// av.v.pageGroups.PageEdit.ButtonAddEmail = tview.NewButton("Добавить e-mail")

	//! "Удалить e-mail"
	// av.v.pageGroups.PageEdit.ButtonDeleteEmail = tview.NewButton("Удалить e-mail")

	//! "Назад"
	av.v.pageGroups.pageEditEmails.buttonEхit.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
	})

	//! "EditEMailsButtonsGrid"
	av.v.pageGroups.pageEditEmails.grid.
		SetRows(3, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageEditEmails.formAddEmail.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageEditEmails.buttonAdd, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageEditEmails.buttonDelete, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageEditEmails.buttonEхit, 3, 0, 1, 1, 0, 0, true)

	//! "emails list"
	av.v.pageGroups.listEmails.ShowSecondaryText(false)
	av.v.pageGroups.listEmails.SetBorder(true)
	av.v.pageGroups.listEmails.SetHighlightFullLine(true)
	av.v.pageGroups.listEmails.SetTitle(" Список допущенных e-mail ")

	//! InputCapture edit emails page
	av.v.pageGroups.pageEditEmails.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.listEmails:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEditEmails.formAddEmail.inputEmail)
			case av.v.pageGroups.pageEditEmails.formAddEmail.inputEmail:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEditEmails.buttonAdd)
			case av.v.pageGroups.pageEditEmails.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEditEmails.buttonDelete)
			case av.v.pageGroups.pageEditEmails.buttonDelete:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEditEmails.buttonEхit)
			case av.v.pageGroups.pageEditEmails.buttonEхit:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.listEmails)
			}
			return nil
		}
		return event
	}

	//! Main grid
	av.v.pageGroups.gridMain.
		SetRows(0).
		SetColumns(0, 30, 20, 30, 0).
		SetGap(1, 1).
		AddItem(av.v.pageGroups.listGroups, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pages, 0, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.listEmails, 0, 3, 1, 1, 0, 0, true)

	//! adding pages
	av.v.pageGroups.pages.AddPage("select_page", av.v.pageGroups.pageSelect.grid, true, true)
	av.v.pageGroups.pages.AddPage("add_group_page", av.v.pageGroups.pageAddGroup.grid, true, true)
	av.v.pageGroups.pages.AddPage("edit_emails_page", av.v.pageGroups.pageEditEmails.grid, true, true)
}
