package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataViewNote struct {
	textviewNoteL *tview.TextView
	textviewDescL *tview.TextView
	textviewNote  *tview.TextView
	textviewDesc  *tview.TextView
	gridData      *tview.Grid
	inputCapture  func(event *tcell.EventKey) *tcell.EventKey
	page          *tview.Pages
}

func newPageDataViewNote() *pageDataViewNote {
	return &pageDataViewNote{
		textviewNoteL: tview.NewTextView(),
		textviewDescL: tview.NewTextView(),
		textviewNote:  tview.NewTextView(),
		textviewDesc:  tview.NewTextView(),
		gridData:      tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataViewNote() {
	//! note
	av.v.pageData.pageDataViewNote.textviewNoteL.SetText("Заметка:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewNote.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)

	av.v.pageData.pageDataViewNote.textviewNote.SetText("В Telegram-канале самой компании в 12:05 мск сообщили, что «Крок» продолжает свою работу в штатном режиме. Все бизнес-процессы, включая поддержку клиентов, функционируют в рамках установленных регламентов и осуществляются без перебоев. По данным из открытых источников, в 2021 году «Крок» занимала девятое место по выручке среди всех российских IT-компаний. Она специализируется на IT-услугах в области системной интеграции, Big Data, блокчейне, искусственном интеллекте, машинном обучении и других.")
	av.v.pageData.pageDataViewNote.textviewDesc.SetText("Силовики приехали c обысками в офис одной из крупнейших IT-компаний «Крок» в Москве. Об этом ТАСС сообщили в оперативных службах.")

	//! view note page grid

	av.v.pageData.pageDataViewNote.gridData.
		SetBorders(true).
		SetRows(0, 4).
		SetColumns(9, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataViewNote.textviewNoteL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewNote.textviewDescL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewNote.textviewNote, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewNote.textviewDesc, 1, 1, 1, 1, 0, 0, true)
		// SetBorder(true).
		// SetTitle(" Данные ")

	av.v.pageData.pageDataViewNote.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
}
