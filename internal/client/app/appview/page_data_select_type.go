package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataSelectType struct {
	buttonNote   *tview.Button
	buttonPass   *tview.Button
	buttonCard   *tview.Button
	buttonFile   *tview.Button
	buttonCancel *tview.Button
	grid         *tview.Grid
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
	page         *tview.Pages
}

func newPageDataSelectType() *pageDataSelectType {
	return &pageDataSelectType{
		buttonNote:   tview.NewButton("Заметка"),
		buttonPass:   tview.NewButton("Пароль"),
		buttonCard:   tview.NewButton("Карта"),
		buttonFile:   tview.NewButton("Файл"),
		buttonCancel: tview.NewButton("Отмена"),
		grid:         tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataSelectType() {
	//! "Заметка"
	av.v.pageData.pageDataSelectType.buttonNote.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_add_note_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataAddNote.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.textareaNote)
	})

	//! "Пароль"
	av.v.pageData.pageDataSelectType.buttonPass.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_add_pass_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataAddPass.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.textareaLogin)
	})

	//! "Карта"
	av.v.pageData.pageDataSelectType.buttonCard.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_add_card_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataAddCard.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaNumber)
	})

	//! "Файл"
	av.v.pageData.pageDataSelectType.buttonFile.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_add_file_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataAddFile.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddFile.treeview)
	})

	//! "Отмена"
	av.v.pageData.pageDataSelectType.buttonCancel.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_view_note_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataViewNote.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
	})

	//! grid
	av.v.pageData.pageDataSelectType.grid.
		SetRows(1, 1, 1, 1, 1, 1).
		SetColumns(5, 20, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataSelectType.buttonNote, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataSelectType.buttonPass, 2, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataSelectType.buttonCard, 3, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataSelectType.buttonFile, 4, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataSelectType.buttonCancel, 5, 1, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataSelectType.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.pageDataSelectType.buttonNote:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataSelectType.buttonPass)
			case av.v.pageData.pageDataSelectType.buttonPass:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataSelectType.buttonCard)
			case av.v.pageData.pageDataSelectType.buttonCard:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataSelectType.buttonFile)
			case av.v.pageData.pageDataSelectType.buttonFile:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataSelectType.buttonCancel)
			case av.v.pageData.pageDataSelectType.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataSelectType.buttonNote)
			}
			return nil
		}
		return event
	}
}
