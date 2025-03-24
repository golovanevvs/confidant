package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataAddCard struct {
	textviewNumberL *tview.TextView
	textviewDateL   *tview.TextView
	textviewNameL   *tview.TextView
	textviewCVC2L   *tview.TextView
	textviewPINL    *tview.TextView
	textviewBankL   *tview.TextView
	textviewDescL   *tview.TextView
	textviewTitleL  *tview.TextView
	textareaNumber  *tview.TextArea
	textareaDate    *tview.TextArea
	textareaName    *tview.TextArea
	textareaCVC2    *tview.TextArea
	textareaPIN     *tview.TextArea
	textareaBank    *tview.TextArea
	textareaDesc    *tview.TextArea
	textareaTitle   *tview.TextArea
	buttonAdd       *tview.Button
	buttonCancel    *tview.Button
	gridData        *tview.Grid
	gridButtons     *tview.Grid
	grid            *tview.Grid
	inputCapture    func(event *tcell.EventKey) *tcell.EventKey
	page            *tview.Pages
}

func newPageDataAddCard() *pageDataAddCard {
	return &pageDataAddCard{
		textviewNumberL: tview.NewTextView(),
		textviewDateL:   tview.NewTextView(),
		textviewNameL:   tview.NewTextView(),
		textviewCVC2L:   tview.NewTextView(),
		textviewPINL:    tview.NewTextView(),
		textviewBankL:   tview.NewTextView(),
		textviewDescL:   tview.NewTextView(),
		textviewTitleL:  tview.NewTextView(),
		textareaNumber:  tview.NewTextArea(),
		textareaDate:    tview.NewTextArea(),
		textareaName:    tview.NewTextArea(),
		textareaCVC2:    tview.NewTextArea(),
		textareaPIN:     tview.NewTextArea(),
		textareaBank:    tview.NewTextArea(),
		textareaDesc:    tview.NewTextArea(),
		textareaTitle:   tview.NewTextArea(),
		buttonAdd:       tview.NewButton("Добавить"),
		buttonCancel:    tview.NewButton("Отмена"),
		gridData:        tview.NewGrid(),
		gridButtons:     tview.NewGrid(),
		grid:            tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataAddCard() {
	//! label names
	av.v.pageData.pageDataAddCard.textviewNumberL.SetText("Номер:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewDateL.SetText("Годна до:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewNameL.SetText("Имя:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewCVC2L.SetText("CVC2:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewPINL.SetText("PIN:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewBankL.SetText("Банк:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddCard.textviewTitleL.SetText("Название:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataAddCard.gridData.
		SetBorders(true).
		SetRows(1, 1, 1, 1, 0, 1).
		SetColumns(9, 15, 4, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataAddCard.textviewNumberL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewDateL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewNameL, 1, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewCVC2L, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewPINL, 2, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewBankL, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewDescL, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaNumber, 0, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaDate, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaName, 1, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaCVC2, 2, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaPIN, 2, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaBank, 3, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaDesc, 4, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textviewTitleL, 5, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.textareaTitle, 5, 1, 1, 3, 0, 0, true)

	//! Добавить

	//! Отмена
	av.v.pageData.pageDataAddCard.buttonCancel.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_view_card_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataViewCard.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
	})

	//! buttons grid
	av.v.pageData.pageDataAddCard.gridButtons.
		SetBorders(false).
		SetRows(1).
		SetColumns(10, 10).
		SetGap(0, 1).
		AddItem(av.v.pageData.pageDataAddCard.buttonAdd, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.buttonCancel, 0, 1, 1, 1, 0, 0, true)

	//! grid
	av.v.pageData.pageDataAddCard.grid.
		SetBorders(false).
		SetRows(0, 1).
		SetColumns(0).
		AddItem(av.v.pageData.pageDataAddCard.gridData, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddCard.gridButtons, 1, 0, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataAddCard.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.pageDataAddCard.textareaNumber:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaDate)
			case av.v.pageData.pageDataAddCard.textareaDate:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaName)
			case av.v.pageData.pageDataAddCard.textareaName:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaCVC2)
			case av.v.pageData.pageDataAddCard.textareaCVC2:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaPIN)
			case av.v.pageData.pageDataAddCard.textareaPIN:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaBank)
			case av.v.pageData.pageDataAddCard.textareaBank:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaDesc)
			case av.v.pageData.pageDataAddCard.textareaDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaTitle)
			case av.v.pageData.pageDataAddCard.textareaTitle:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.buttonAdd)
			case av.v.pageData.pageDataAddCard.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.buttonCancel)
			case av.v.pageData.pageDataAddCard.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddCard.textareaNumber)
			}
			return nil
		}
		return event
	}

}
