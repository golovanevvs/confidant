package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataViewPass struct {
	textviewLoginL *tview.TextView
	textviewPassL  *tview.TextView
	textviewDescL  *tview.TextView
	textviewLogin  *tview.TextView
	textviewPass   *tview.TextView
	textviewDesc   *tview.TextView
	gridData       *tview.Grid
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	page           *tview.Pages
}

func newPageDataViewPass() *pageDataViewPass {
	return &pageDataViewPass{
		textviewLoginL: tview.NewTextView(),
		textviewPassL:  tview.NewTextView(),
		textviewDescL:  tview.NewTextView(),
		textviewLogin:  tview.NewTextView(),
		textviewPass:   tview.NewTextView(),
		textviewDesc:   tview.NewTextView(),
		gridData:       tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataViewPass() {
	av.v.pageData.pageDataViewPass.textviewLoginL.SetText("Логин:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewPass.textviewPassL.SetText("Пароль:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewPass.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataViewPass.gridData.
		SetBorders(true).
		SetRows(1, 1, 0).
		SetColumns(9, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataViewPass.textviewLoginL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewPass.textviewPassL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewPass.textviewDescL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewPass.textviewLogin, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewPass.textviewPass, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewPass.textviewDesc, 2, 1, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataViewPass.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewPass.textviewDesc)
			case av.v.pageData.pageDataViewPass.textviewDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}
}
