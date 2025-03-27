package appview

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/rivo/tview"
)

type pageDataAddNote struct {
	textviewNoteL  *tview.TextView
	textviewDescL  *tview.TextView
	textviewTitleL *tview.TextView
	textareaNote   *tview.TextArea
	textareaDesc   *tview.TextArea
	textareaTitle  *tview.TextArea
	buttonAdd      *tview.Button
	buttonCancel   *tview.Button
	gridData       *tview.Grid
	gridButtons    *tview.Grid
	grid           *tview.Grid
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	page           *tview.Pages
}

func newPageDataAddNote() *pageDataAddNote {
	return &pageDataAddNote{
		textviewNoteL:  tview.NewTextView(),
		textviewDescL:  tview.NewTextView(),
		textviewTitleL: tview.NewTextView(),
		textareaNote:   tview.NewTextArea(),
		textareaDesc:   tview.NewTextArea(),
		textareaTitle:  tview.NewTextArea(),
		buttonAdd:      tview.NewButton("Добавить"),
		buttonCancel:   tview.NewButton("Отмена"),
		gridData:       tview.NewGrid(),
		gridButtons:    tview.NewGrid(),
		grid:           tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataAddNote() {

	//! label names
	av.v.pageData.pageDataAddNote.textviewNoteL.SetText("Заметка:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddNote.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddNote.textviewTitleL.SetText("Название:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataAddNote.gridData.
		SetBorders(true).
		SetRows(0, 4, 1).
		SetColumns(9, 0).
		AddItem(av.v.pageData.pageDataAddNote.textviewNoteL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.textviewDescL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.textviewTitleL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.textareaNote, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.textareaDesc, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.textareaTitle, 2, 1, 1, 1, 0, 0, true)

	//! Добавить
	av.v.pageData.pageDataAddNote.buttonAdd.SetSelectedFunc(func() {
		title := av.v.pageData.pageDataAddNote.textareaTitle.GetText()
		desc := av.v.pageData.pageDataAddNote.textareaDesc.GetText()
		note := av.v.pageData.pageDataAddNote.textareaNote.GetText()
		data := model.NoteDec{
			Desc:  desc,
			Note:  note,
			Title: title,
		}
		err := av.sv.AddNote(context.Background(), data, av.account.ID, av.groupID)
		if err != nil {
			av.v.pageMain.messageBoxL.SetText(err.Error())
		} else {
			av.aPageDataSwitch()
		}
	})

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

	//! grid
	av.v.pageData.pageDataAddNote.grid.
		SetBorders(false).
		SetRows(0, 1).
		SetColumns(0).
		AddItem(av.v.pageData.pageDataAddNote.gridData, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddNote.gridButtons, 1, 0, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataAddNote.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.pageDataAddNote.textareaNote:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.textareaDesc)
			case av.v.pageData.pageDataAddNote.textareaDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.textareaTitle)
			case av.v.pageData.pageDataAddNote.textareaTitle:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.buttonAdd)
			case av.v.pageData.pageDataAddNote.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.buttonCancel)
			case av.v.pageData.pageDataAddNote.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddNote.textareaNote)
			}
			return nil
		}
		return event
	}
}
