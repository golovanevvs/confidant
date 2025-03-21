package appview

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// type pageViewPassword struct {
// }

// type pageViewCard struct {
// }

// type pageViewFile struct {
// }

type pageData struct {
	listTitles       *tview.List
	buttonAdd        *tview.Button
	buttonEdit       *tview.Button
	buttonDelete     *tview.Button
	buttonBack       *tview.Button
	buttonExit       *tview.Button
	gridButtons      *tview.Grid
	gridMain         *tview.Grid
	pageDataViewNote *pageDataViewNote
	pageDataAddNote  *pageDataAddNote
	pages            *tview.Pages
	inputCapture     func(event *tcell.EventKey) *tcell.EventKey
	page             *tview.Pages
}

func newPageData() *pageData {
	return &pageData{
		listTitles:       tview.NewList(),
		buttonAdd:        tview.NewButton("Добавить"),
		buttonEdit:       tview.NewButton("Изменить"),
		buttonDelete:     tview.NewButton("Удалить"),
		buttonBack:       tview.NewButton("Назад"),
		buttonExit:       tview.NewButton("Выход"),
		gridButtons:      tview.NewGrid(),
		gridMain:         tview.NewGrid(),
		pageDataViewNote: newPageDataViewNote(),
		pageDataAddNote:  newPageDataAddNote(),
		pages:            tview.NewPages(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vData() {
	//! titles list
	av.v.pageData.listTitles.ShowSecondaryText(false)
	av.v.pageData.listTitles.SetBorder(true)
	av.v.pageData.listTitles.SetHighlightFullLine(true)
	for i := 0; i < 10; i++ {
		av.v.pageData.listTitles.AddItem(fmt.Sprintf("Title %d", i), "", 0, nil)
	}

	av.v.pageData.listTitles.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	av.v.pageData.listTitles.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText)
	})

	//! "Добавить"
	av.v.pageData.buttonAdd.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_add_note")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataAddNote.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.textareaNote)
	})

	//! "Назад"
	av.v.pageData.buttonBack.SetSelectedFunc(func() {
		av.v.pageMain.pages.SwitchToPage("groups_page")
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
	})

	//! "Выход"
	av.v.pageData.buttonExit.SetSelectedFunc(func() {
		av.v.pageApp.app.Stop()
	})

	//! buttons grid
	av.v.pageData.gridButtons.
		SetRows(1, 1, 1, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageData.buttonAdd, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.buttonEdit, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.buttonDelete, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.buttonBack, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.buttonExit, 5, 0, 1, 1, 0, 0, true)

	//! Main grid
	av.v.pageData.gridMain.
		SetRows(0).
		SetColumns(30, 20, 0).
		AddItem(av.v.pageData.listTitles, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.gridButtons, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pages, 0, 2, 1, 1, 0, 0, true)

	av.v.pageData.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.listTitles:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonAdd)
			case av.v.pageData.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonEdit)
			case av.v.pageData.buttonEdit:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonDelete)
			case av.v.pageData.buttonDelete:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonBack)
			case av.v.pageData.buttonBack:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonExit)
			case av.v.pageData.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewNote.textviewNote)
			case av.v.pageData.pageDataViewNote.textviewNote:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewNote.textviewDesc)
			case av.v.pageData.pageDataViewNote.textviewDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}

	//! adding pages
	av.v.pageData.pages.AddPage("data_add_note", av.v.pageData.pageDataAddNote.grid, true, true)
	av.v.pageData.pages.AddPage("data_view_note_page", av.v.pageData.pageDataViewNote.gridData, true, true)
}
