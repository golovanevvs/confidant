package appview

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PageSelect struct {
	GridSelectButtons *tview.Grid
	ButtonSelect      *tview.Button
	ButtonNew         *tview.Button
	ButtonSettings    *tview.Button
	ButtonDelete      *tview.Button
	ButtonLogout      *tview.Button
	ButtonExit        *tview.Button
	Grid              *tview.Grid
	InputCapture      func(event *tcell.EventKey) *tcell.EventKey
	Page              *tview.Pages
}

type PageEdit struct {
	FormAddEmail      FormAdd
	ButtonAddEmail    *tview.Button
	ButtonDeleteEmail *tview.Button
	ButtonEditExit    *tview.Button
	Grid              *tview.Grid
	InputCapture      func(event *tcell.EventKey) *tcell.EventKey
	Page              *tview.Pages
}

type FormAdd struct {
	InputEmail *tview.InputField
	Form       *tview.Form
}

type PageGroups struct {
	ListGroups *tview.List
	ListEmails *tview.List
	GridMain   *tview.Grid
	PageSelect PageSelect
	PageEdit   PageEdit
	PagesSelEd *tview.Pages
}

func (av *AppView) VGroups() {
	//? pages
	// av.v.pageGroups.PagesSelEd = tview.NewPages()

	//? Groups List
	// av.v.pageGroups.ListGroups = tview.NewList()
	av.v.pageGroups.ListGroups.SetBorder(true)
	av.v.pageGroups.ListGroups.SetHighlightFullLine(true)
	av.v.pageGroups.ListGroups.SetTitle(" Список групп ")
	for i := 0; i < 10; i++ {
		av.v.pageGroups.ListGroups.AddItem(fmt.Sprintf("Group %d", i), "", 0, nil)
	}

	av.v.pageGroups.ListGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.MessageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	//? "select_page"

	//? Выбрать группу
	// av.v.pageGroups.PageSelect.ButtonSelect = tview.NewButton("Выбрать группу")

	//? Создать группу
	// av.v.pageGroups.PageSelect.ButtonNew = tview.NewButton("Создать группу")

	//? Настроить группу
	av.v.pageGroups.PageSelect.ButtonSettings.SetSelectedFunc(func() {
		av.v.pageGroups.PagesSelEd.SwitchToPage("edit_page")
		av.v.pageMain.App.SetInputCapture(av.v.pageGroups.PageEdit.InputCapture)
		av.v.pageMain.App.SetFocus(av.v.pageGroups.ListEmails)
		av.v.pageGroups.PageEdit.FormAddEmail.InputEmail.SetText("")
	})

	//? Удалить группу
	// av.v.pageGroups.PageSelect.ButtonDelete = tview.NewButton("Удалить группу")

	//? Выйти из аккаунта
	av.v.pageGroups.PageSelect.ButtonLogout.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.Pages.SwitchToPage("login_page")
		// focus
		av.v.pageMain.App.SetInputCapture(av.v.pageLogin.InputCapture)
		av.v.pageMain.App.SetFocus(av.v.pageLogin.Form.InputEmail)
		// messageBox
		av.v.pageMain.MessageBoxL.Clear()
		av.v.pageMain.MessageBoxR.Clear()
	})

	//? Выход
	av.v.pageGroups.PageSelect.ButtonExit.SetSelectedFunc(func() {
		av.v.pageMain.App.Stop()
	})

	//? MainButtonsGrid
	av.v.pageGroups.PageSelect.Grid.
		SetRows(1, 1, 1, 1, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.PageSelect.ButtonSelect, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageSelect.ButtonNew, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageSelect.ButtonSettings, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageSelect.ButtonDelete, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageSelect.ButtonLogout, 5, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageSelect.ButtonExit, 6, 0, 1, 1, 0, 0, true)

	//? InputCapture select page
	av.v.pageGroups.PageSelect.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageMain.App.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.ListGroups:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageSelect.ButtonSelect)
			case av.v.pageGroups.PageSelect.ButtonSelect:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageSelect.ButtonNew)
			case av.v.pageGroups.PageSelect.ButtonNew:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageSelect.ButtonSettings)
			case av.v.pageGroups.PageSelect.ButtonSettings:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageSelect.ButtonDelete)
			case av.v.pageGroups.PageSelect.ButtonDelete:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageSelect.ButtonLogout)
			case av.v.pageGroups.PageSelect.ButtonLogout:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageSelect.ButtonExit)
			case av.v.pageGroups.PageSelect.ButtonExit:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.ListGroups)
			}
			return nil
		}
		return event
	}

	//? "edit_page"

	//? add form
	av.v.pageGroups.PageEdit.FormAddEmail.Form.SetHorizontal(false)
	av.v.pageGroups.PageEdit.FormAddEmail.Form.AddInputField("", "", 0, nil, nil)
	av.v.pageGroups.PageEdit.FormAddEmail.InputEmail = av.v.pageGroups.PageEdit.FormAddEmail.Form.GetFormItem(0).(*tview.InputField)

	//? Добавить e-mail
	// av.v.pageGroups.PageEdit.ButtonAddEmail = tview.NewButton("Добавить e-mail")

	//? Удалить e-mail
	// av.v.pageGroups.PageEdit.ButtonDeleteEmail = tview.NewButton("Удалить e-mail")

	//? Назад
	av.v.pageGroups.PageEdit.ButtonEditExit.SetSelectedFunc(func() {
		av.v.pageGroups.PagesSelEd.SwitchToPage("select_page")
		av.v.pageMain.App.SetInputCapture(av.v.pageGroups.PageSelect.InputCapture)
		av.v.pageMain.App.SetFocus(av.v.pageGroups.ListGroups)
	})

	//? EditButtonsGrid
	av.v.pageGroups.PageEdit.Grid.
		SetRows(3, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.PageEdit.FormAddEmail.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageEdit.ButtonAddEmail, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageEdit.ButtonDeleteEmail, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PageEdit.ButtonEditExit, 3, 0, 1, 1, 0, 0, true)

	//? emails list
	av.v.pageGroups.ListEmails.SetBorder(true)
	av.v.pageGroups.ListEmails.SetHighlightFullLine(true)
	av.v.pageGroups.ListEmails.SetTitle(" Список допущенных e-mail ")

	//? InputCapture edit page
	av.v.pageGroups.PageEdit.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageMain.App.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.ListEmails:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageEdit.FormAddEmail.InputEmail)
			case av.v.pageGroups.PageEdit.FormAddEmail.InputEmail:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageEdit.ButtonAddEmail)
			case av.v.pageGroups.PageEdit.ButtonAddEmail:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageEdit.ButtonDeleteEmail)
			case av.v.pageGroups.PageEdit.ButtonDeleteEmail:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.PageEdit.ButtonEditExit)
			case av.v.pageGroups.PageEdit.ButtonEditExit:
				av.v.pageMain.App.SetFocus(av.v.pageGroups.ListEmails)
			}
			return nil
		}
		return event
	}

	//? Main grid
	av.v.pageGroups.GridMain.
		SetRows(0).
		SetColumns(0, 30, 20, 30, 0).
		SetGap(1, 1).
		AddItem(av.v.pageGroups.ListGroups, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.PagesSelEd, 0, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.ListEmails, 0, 3, 1, 1, 0, 0, true)

		//? Adding pages
	av.v.pageGroups.PagesSelEd.AddPage("select_page", av.v.pageGroups.PageSelect.Grid, true, true)
	av.v.pageGroups.PagesSelEd.AddPage("edit_page", av.v.pageGroups.PageEdit.Grid, true, true)
}
