package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageGroupsEditEmails struct {
	formGroupsAddEmail *formGroupsAddEmail
	buttonAdd          *tview.Button
	buttonDelete       *tview.Button
	buttonEхit         *tview.Button
	grid               *tview.Grid
	inputCapture       func(event *tcell.EventKey) *tcell.EventKey
	page               *tview.Pages
}

type formGroupsAddEmail struct {
	inputEmail *tview.InputField
	form       *tview.Form
}

func newPageFroupsEditEmails() *pageGroupsEditEmails {
	return &pageGroupsEditEmails{
		formGroupsAddEmail: &formGroupsAddEmail{
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
	}
}

func (av *appView) vGroupsEditEmails() {
	//! add form
	av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.form.SetHorizontal(false)
	av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.form.AddInputField("", "", 0, nil, nil)
	av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.inputEmail = av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.form.GetFormItem(0).(*tview.InputField)

	//! "Добавить e-mail"
	// av.v.pageGroups.PageEdit.ButtonAddEmail = tview.NewButton("Добавить e-mail")

	//! "Удалить e-mail"
	// av.v.pageGroups.PageEdit.ButtonDeleteEmail = tview.NewButton("Удалить e-mail")

	//! "Назад"
	av.v.pageGroups.pageGroupsEditEmails.buttonEхit.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
	})

	//! "EditEMailsButtonsGrid"
	av.v.pageGroups.pageGroupsEditEmails.grid.
		SetRows(3, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsEditEmails.buttonAdd, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsEditEmails.buttonDelete, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsEditEmails.buttonEхit, 3, 0, 1, 1, 0, 0, true)

	//! "emails list"
	av.v.pageGroups.listEmails.ShowSecondaryText(false)
	av.v.pageGroups.listEmails.SetBorder(true)
	av.v.pageGroups.listEmails.SetHighlightFullLine(true)
	av.v.pageGroups.listEmails.SetTitle(" Список допущенных e-mail ")

	//! InputCapture edit emails page
	av.v.pageGroups.pageGroupsEditEmails.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.listEmails:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.inputEmail)
			case av.v.pageGroups.pageGroupsEditEmails.formGroupsAddEmail.inputEmail:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsEditEmails.buttonAdd)
			case av.v.pageGroups.pageGroupsEditEmails.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsEditEmails.buttonDelete)
			case av.v.pageGroups.pageGroupsEditEmails.buttonDelete:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsEditEmails.buttonEхit)
			case av.v.pageGroups.pageGroupsEditEmails.buttonEхit:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.listEmails)
			}
			return nil
		}
		return event
	}
}
