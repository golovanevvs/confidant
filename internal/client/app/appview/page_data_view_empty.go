package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataViewEmpty struct {
	textviewEmpty *tview.TextView
	gridData      *tview.Grid
	inputCapture  func(event *tcell.EventKey) *tcell.EventKey
	page          *tview.Pages
}

func newPageDataViewEmpty() *pageDataViewEmpty {
	return &pageDataViewEmpty{
		textviewEmpty: tview.NewTextView(),
		gridData:      tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataViewEmpty() {
	av.v.pageData.pageDataViewEmpty.textviewEmpty.SetText("Пусто").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataViewEmpty.gridData.
		SetRows(0, 1, 0).
		SetColumns(0, 5, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataViewEmpty.textviewEmpty, 1, 1, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataViewEmpty.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.listTitles:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonAdd)
			case av.v.pageData.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonBack)
			// case av.v.pageData.buttonEdit:
			// 	av.v.pageApp.app.SetFocus(av.v.pageData.buttonDelete)
			// case av.v.pageData.buttonDelete:
			// 	av.v.pageApp.app.SetFocus(av.v.pageData.buttonBack)
			case av.v.pageData.buttonBack:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonExit)
			case av.v.pageData.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}
}

func (av *appView) aPageDataViewEmptySwitch() {
	av.v.pageMain.pages.SwitchToPage("data_page")
	av.v.pageData.pages.SwitchToPage("data_view_empty_page")
	av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataViewEmpty.inputCapture)
	av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
}
