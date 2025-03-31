package appview

import (
	"context"

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

	//! data grid

	av.v.pageData.pageDataViewNote.gridData.
		SetBorders(true).
		SetRows(0, 4).
		SetColumns(9, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataViewNote.textviewNoteL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewNote.textviewDescL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewNote.textviewNote, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewNote.textviewDesc, 1, 1, 1, 1, 0, 0, true)

	//! inputCapture
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

func (av *appView) vPageDataViewNoteUpdate() {
	av.vClearMessages()

	data, err := av.sv.GetNote(context.Background(), av.dataID)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(err.Error())
	} else {
		av.v.pageData.pageDataViewNote.textviewNote.SetText(data.Note)
		av.v.pageData.pageDataViewNote.textviewDesc.SetText(data.Desc)
		av.v.pageMain.pages.SwitchToPage("data_page")
		av.v.pageData.pages.SwitchToPage("data_view_note_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
	}
}
