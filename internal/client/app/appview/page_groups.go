package appview

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageSelect struct {
	gridSelectButtons *tview.Grid
	buttonSelect      *tview.Button
	buttonNew         *tview.Button
	buttonSettings    *tview.Button
	buttonDelete      *tview.Button
	buttonLogout      *tview.Button
	buttonExit        *tview.Button
	grid              *tview.Grid
	inputCapture      func(event *tcell.EventKey) *tcell.EventKey
	Page              *tview.Pages
}

type pageEdit struct {
	formAddEmail      formAdd
	buttonAddEmail    *tview.Button
	buttonDeleteEmail *tview.Button
	buttonEditExit    *tview.Button
	grid              *tview.Grid
	InputCapture      func(event *tcell.EventKey) *tcell.EventKey
	page              *tview.Pages
}

type formAdd struct {
	inputEmail *tview.InputField
	form       *tview.Form
}

type pageGroups struct {
	listGroups *tview.List
	listEmails *tview.List
	gridMain   *tview.Grid
	pageSelect pageSelect
	pageEdit   pageEdit
	pagesSelEd *tview.Pages
}

func newPageGroups() *pageGroups {
	return &pageGroups{
		listGroups: tview.NewList(),
		listEmails: tview.NewList(),
		gridMain:   tview.NewGrid(),
		pageSelect: pageSelect{
			gridSelectButtons: tview.NewGrid(),
			buttonSelect:      tview.NewButton("Выбрать группу"),
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
		pageEdit: pageEdit{
			formAddEmail: formAdd{
				inputEmail: tview.NewInputField(),
				form:       tview.NewForm(),
			},
			buttonAddEmail:    tview.NewButton("Добавить e-mail"),
			buttonDeleteEmail: tview.NewButton("Удалить e-mail"),
			buttonEditExit:    tview.NewButton("Назад"),
			grid:              tview.NewGrid(),
			InputCapture: func(event *tcell.EventKey) *tcell.EventKey {
				return event
			},
			page: tview.NewPages(),
		},
		pagesSelEd: tview.NewPages(),
	}
}

func (av *appView) vGroups() {
	//! Groups List
	av.v.pageGroups.listGroups.SetBorder(true)
	av.v.pageGroups.listGroups.SetHighlightFullLine(true)
	av.v.pageGroups.listGroups.SetTitle(" Список групп ")
	for i := 0; i < 10; i++ {
		av.v.pageGroups.listGroups.AddItem(fmt.Sprintf("Group %d", i), "", 0, nil)
	}

	av.v.pageGroups.listGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	//! "select_page"

	//! "Выбрать группу"
	// av.v.pageGroups.PageSelect.ButtonSelect = tview.NewButton("Выбрать группу")

	//! "Создать группу"
	// av.v.pageGroups.PageSelect.ButtonNew = tview.NewButton("Создать группу")

	//! "Настроить группу"
	av.v.pageGroups.pageSelect.buttonSettings.SetSelectedFunc(func() {
		av.v.pageGroups.pagesSelEd.SwitchToPage("edit_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageEdit.InputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listEmails)
		av.v.pageGroups.pageEdit.formAddEmail.inputEmail.SetText("")
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
		// messageBox
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
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
		AddItem(av.v.pageGroups.pageSelect.buttonSelect, 1, 0, 1, 1, 0, 0, true).
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
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageSelect.buttonSelect)
			case av.v.pageGroups.pageSelect.buttonSelect:
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

	//! "edit_page"

	//! add form
	av.v.pageGroups.pageEdit.formAddEmail.form.SetHorizontal(false)
	av.v.pageGroups.pageEdit.formAddEmail.form.AddInputField("", "", 0, nil, nil)
	av.v.pageGroups.pageEdit.formAddEmail.inputEmail = av.v.pageGroups.pageEdit.formAddEmail.form.GetFormItem(0).(*tview.InputField)

	//! "Добавить e-mail"
	// av.v.pageGroups.PageEdit.ButtonAddEmail = tview.NewButton("Добавить e-mail")

	//! "Удалить e-mail"
	// av.v.pageGroups.PageEdit.ButtonDeleteEmail = tview.NewButton("Удалить e-mail")

	//! "Назад"
	av.v.pageGroups.pageEdit.buttonEditExit.SetSelectedFunc(func() {
		av.v.pageGroups.pagesSelEd.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
	})

	//! "EditButtonsGrid"
	av.v.pageGroups.pageEdit.grid.
		SetRows(3, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageEdit.formAddEmail.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageEdit.buttonAddEmail, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageEdit.buttonDeleteEmail, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageEdit.buttonEditExit, 3, 0, 1, 1, 0, 0, true)

	//! "emails list"
	av.v.pageGroups.listEmails.SetBorder(true)
	av.v.pageGroups.listEmails.SetHighlightFullLine(true)
	av.v.pageGroups.listEmails.SetTitle(" Список допущенных e-mail ")

	//! InputCapture edit page
	av.v.pageGroups.pageEdit.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.listEmails:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEdit.formAddEmail.inputEmail)
			case av.v.pageGroups.pageEdit.formAddEmail.inputEmail:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEdit.buttonAddEmail)
			case av.v.pageGroups.pageEdit.buttonAddEmail:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEdit.buttonDeleteEmail)
			case av.v.pageGroups.pageEdit.buttonDeleteEmail:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageEdit.buttonEditExit)
			case av.v.pageGroups.pageEdit.buttonEditExit:
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
		AddItem(av.v.pageGroups.pagesSelEd, 0, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.listEmails, 0, 3, 1, 1, 0, 0, true)

	//! adding pages
	av.v.pageGroups.pagesSelEd.AddPage("select_page", av.v.pageGroups.pageSelect.grid, true, true)
	av.v.pageGroups.pagesSelEd.AddPage("edit_page", av.v.pageGroups.pageEdit.grid, true, true)
}
