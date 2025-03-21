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
	textareaNumber  *tview.TextArea
	textareaDate    *tview.TextArea
	textareaName    *tview.TextArea
	textareaCVC2    *tview.TextArea
	textareaPIN     *tview.TextArea
	textareaBank    *tview.TextArea
	textareaDesc    *tview.TextArea
	gridData        *tview.Grid
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
		textareaNumber:  tview.NewTextArea(),
		textareaDate:    tview.NewTextArea(),
		textareaName:    tview.NewTextArea(),
		textareaCVC2:    tview.NewTextArea(),
		textareaPIN:     tview.NewTextArea(),
		textareaBank:    tview.NewTextArea(),
		textareaDesc:    tview.NewTextArea(),
		gridData:        tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataAddCard() {
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

	//! data grid
	av.v.pageData.pageDataAddCard.gridData.
		SetBorders(true).
		SetRows(1, 1, 1, 1, 0).
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
		AddItem(av.v.pageData.pageDataAddCard.textareaDesc, 4, 1, 1, 3, 0, 0, true)
}
