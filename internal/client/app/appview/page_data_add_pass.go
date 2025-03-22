package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataAddPass struct {
	textviewLoginL *tview.TextView
	textviewPassL  *tview.TextView
	textviewDescL  *tview.TextView
	textareaLogin  *tview.TextArea
	textareaPass   *tview.TextArea
	textareaDesc   *tview.TextArea
	gridData       *tview.Grid
	gridButtons    *tview.Grid
	grid           *tview.Grid
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	page           *tview.Pages
}

func newPageDataAddPass() *pageDataAddPass {
	return &pageDataAddPass{
		textviewLoginL: tview.NewTextView(),
		textviewPassL:  tview.NewTextView(),
		textviewDescL:  tview.NewTextView(),
		textareaLogin:  tview.NewTextArea(),
		textareaPass:   tview.NewTextArea(),
		textareaDesc:   tview.NewTextArea(),
		gridData:       tview.NewGrid(),
		gridButtons:    tview.NewGrid(),
		grid:           tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataAddPass() {
	//! label names
	av.v.pageData.pageDataAddPass.textviewLoginL.SetText("Логин:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddPass.textviewPassL.SetText("Пароль:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddPass.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataAddPass.gridData.
		SetBorders(true).
		SetRows(1, 1, 0).
		SetColumns(9, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataAddPass.textviewLoginL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textviewPassL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textviewDescL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaLogin, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaPass, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaDesc, 2, 1, 1, 1, 0, 0, true)

	//! Добавить

	//! Отмена
	av.v.pageData.pageDataAddNote.buttonCancel.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_view_note_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataViewNote.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
	})

	//! buttons grid
	av.v.pageData.pageDataAddNote.gridButtons.
		SetBorders(false).
		SetRows(1).
		SetColumns(10, 10).
		SetGap(0, 1).
		AddItem(av.v.pageData.pageDataAddNote.buttonAdd, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.buttonCancel, 0, 1, 1, 1, 0, 0, true)

}
