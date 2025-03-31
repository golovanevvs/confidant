package appview

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataViewCard struct {
	textviewNumberL *tview.TextView
	textviewDateL   *tview.TextView
	textviewNameL   *tview.TextView
	textviewCVC2L   *tview.TextView
	textviewPINL    *tview.TextView
	textviewBankL   *tview.TextView
	textviewDescL   *tview.TextView
	textviewNumber  *tview.TextView
	textviewDate    *tview.TextView
	textviewName    *tview.TextView
	textviewCVC2    *tview.TextView
	textviewPIN     *tview.TextView
	textviewBank    *tview.TextView
	textviewDesc    *tview.TextView
	gridData        *tview.Grid
	inputCapture    func(event *tcell.EventKey) *tcell.EventKey
	page            *tview.Pages
}

func newPageDataViewCard() *pageDataViewCard {
	return &pageDataViewCard{
		textviewNumberL: tview.NewTextView(),
		textviewDateL:   tview.NewTextView(),
		textviewNameL:   tview.NewTextView(),
		textviewCVC2L:   tview.NewTextView(),
		textviewPINL:    tview.NewTextView(),
		textviewBankL:   tview.NewTextView(),
		textviewDescL:   tview.NewTextView(),
		textviewNumber:  tview.NewTextView(),
		textviewDate:    tview.NewTextView(),
		textviewName:    tview.NewTextView(),
		textviewCVC2:    tview.NewTextView(),
		textviewPIN:     tview.NewTextView(),
		textviewBank:    tview.NewTextView(),
		textviewDesc:    tview.NewTextView(),
		gridData:        tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataViewCard() {
	av.v.pageData.pageDataViewCard.textviewNumberL.SetText("Номер:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewCard.textviewDateL.SetText("Годна до:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewCard.textviewNameL.SetText("Имя:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewCard.textviewCVC2L.SetText("CVC2:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewCard.textviewPINL.SetText("PIN:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewCard.textviewBankL.SetText("Банк:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewCard.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataViewCard.gridData.
		SetBorders(true).
		SetRows(1, 1, 1, 1, 0).
		SetColumns(9, 15, 4, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataViewCard.textviewNumberL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewDateL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewNameL, 1, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewCVC2L, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewPINL, 2, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewBankL, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewDescL, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewNumber, 0, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewDate, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewName, 1, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewCVC2, 2, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewPIN, 2, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewBank, 3, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewCard.textviewDesc, 4, 1, 1, 3, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataViewCard.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewCard.textviewDesc)
			case av.v.pageData.pageDataViewCard.textviewDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}

}

func (av *appView) vPageDataViewCardUpdate() {
	av.vClearMessages()

	data, err := av.sv.GetCard(context.Background(), av.dataID)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(err.Error())
	} else {
		av.v.pageData.pageDataViewCard.textviewDesc.SetText(data.Desc)
		av.v.pageData.pageDataViewCard.textviewNumber.SetText(data.Number)
		av.v.pageData.pageDataViewCard.textviewDate.SetText(data.Date)
		av.v.pageData.pageDataViewCard.textviewName.SetText(data.Name)
		av.v.pageData.pageDataViewCard.textviewCVC2.SetText(data.CVC2)
		av.v.pageData.pageDataViewCard.textviewPIN.SetText(data.PIN)
		av.v.pageData.pageDataViewCard.textviewBank.SetText(data.Bank)
		av.v.pageMain.pages.SwitchToPage("data_page")
		av.v.pageData.pages.SwitchToPage("data_view_card_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
	}
}
