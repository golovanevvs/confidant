package appview

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageViewNote struct {
	textviewNoteL *tview.TextView
	textviewDescL *tview.TextView
	textviewNote  *tview.TextView
	textviewDesc  *tview.TextView
	gridData      *tview.Grid
	inputCapture  func(event *tcell.EventKey) *tcell.EventKey
	page          *tview.Pages
}

type pageAddNote struct {
	textviewNoteL  *tview.TextView
	textviewDescL  *tview.TextView
	textviewTitleL *tview.TextView
	textareaNote   *tview.TextArea
	textareaDesc   *tview.TextArea
	textareaTitle  *tview.TextArea
	buttonAdd      *tview.Button
	buttonCancel   *tview.Button
	gridButtons    *tview.Grid
	gridData       *tview.Grid
	grid           *tview.Grid
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	page           *tview.Pages
}

// type pageViewPassword struct {
// }

// type pageViewCard struct {
// }

// type pageViewFile struct {
// }

type pageData struct {
	listTitles   *tview.List
	buttonAdd    *tview.Button
	buttonEdit   *tview.Button
	buttonDelete *tview.Button
	buttonBack   *tview.Button
	buttonExit   *tview.Button
	gridButtons  *tview.Grid
	gridMain     *tview.Grid
	pageViewNote pageViewNote
	pageAddNote  pageAddNote
	pages        *tview.Pages
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
	page         *tview.Pages
}

func newPageData() *pageData {
	return &pageData{
		listTitles:   tview.NewList(),
		buttonAdd:    tview.NewButton("Добавить"),
		buttonEdit:   tview.NewButton("Изменить"),
		buttonDelete: tview.NewButton("Удалить"),
		buttonBack:   tview.NewButton("Назад"),
		buttonExit:   tview.NewButton("Выход"),
		gridButtons:  tview.NewGrid(),
		gridMain:     tview.NewGrid(),
		pageViewNote: pageViewNote{
			textviewNoteL: tview.NewTextView(),
			textviewDescL: tview.NewTextView(),
			textviewNote:  tview.NewTextView(),
			textviewDesc:  tview.NewTextView(),
			gridData:      tview.NewGrid(),
			inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
				return event
			},
			page: tview.NewPages(),
		},
		pageAddNote: pageAddNote{
			textviewNoteL:  tview.NewTextView(),
			textviewDescL:  tview.NewTextView(),
			textviewTitleL: tview.NewTextView(),
			textareaNote:   tview.NewTextArea(),
			textareaDesc:   tview.NewTextArea(),
			textareaTitle:  tview.NewTextArea(),
			buttonAdd:      tview.NewButton("Добавить"),
			buttonCancel:   tview.NewButton("Отмена"),
			gridButtons:    tview.NewGrid(),
			gridData:       tview.NewGrid(),
			grid:           tview.NewGrid(),
			inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
				return event
			},
			page: tview.NewPages(),
		},
		pages: tview.NewPages(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vData() {

	//! MAIN -------------------------------------------------------------------

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
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageAddNote.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageAddNote.textareaNote)
	})

	//! "Назад"
	av.v.pageData.buttonBack.SetSelectedFunc(func() {
		av.v.pageMain.pages.SwitchToPage("groups_page")
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageSelect.inputCapture)
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

	//! VIEW NOTE PAGE -------------------------------------------------------------

	//! note
	av.v.pageData.pageViewNote.textviewNoteL.SetText("Заметка:")
	av.v.pageData.pageViewNote.textviewDescL.SetText("Описание:")
	av.v.pageData.pageViewNote.textviewNote.SetText("В Telegram-канале самой компании в 12:05 мск сообщили, что «Крок» продолжает свою работу в штатном режиме. Все бизнес-процессы, включая поддержку клиентов, функционируют в рамках установленных регламентов и осуществляются без перебоев. По данным из открытых источников, в 2021 году «Крок» занимала девятое место по выручке среди всех российских IT-компаний. Она специализируется на IT-услугах в области системной интеграции, Big Data, блокчейне, искусственном интеллекте, машинном обучении и других.")
	av.v.pageData.pageViewNote.textviewDesc.SetText("Силовики приехали c обысками в офис одной из крупнейших IT-компаний «Крок» в Москве. Об этом ТАСС сообщили в оперативных службах.")

	//! view note page grid

	av.v.pageData.pageViewNote.gridData.
		SetBorders(true).
		SetRows(0, 4).
		SetColumns(9, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageViewNote.textviewNoteL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageViewNote.textviewDescL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageViewNote.textviewNote, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageViewNote.textviewDesc, 1, 1, 1, 1, 0, 0, true)
		// SetBorder(true).
		// SetTitle(" Данные ")

	av.v.pageData.pageViewNote.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
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
				av.v.pageApp.app.SetFocus(av.v.pageData.pageViewNote.textviewNote)
			case av.v.pageData.pageViewNote.textviewNote:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageViewNote.textviewDesc)
			case av.v.pageData.pageViewNote.textviewDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}

	//! ADD NOTE PAGE ---------------------------------------------------------------

	//! note label

	av.v.pageData.pageAddNote.textviewNoteL.SetText("Заметка:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageAddNote.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageAddNote.textviewTitleL.SetText("Название:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! add note page data grid

	av.v.pageData.pageAddNote.gridData.
		SetBorders(true).
		SetRows(0, 4, 1).
		SetColumns(9, 0).
		AddItem(av.v.pageData.pageAddNote.textviewNoteL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.textviewDescL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.textviewTitleL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.textareaNote, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.textareaDesc, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.textareaTitle, 2, 1, 1, 1, 0, 0, true)

	//! add note page buttons grid

	av.v.pageData.pageAddNote.gridButtons.
		SetBorders(false).
		SetRows(1).
		SetColumns(10, 10).
		SetGap(0, 1).
		AddItem(av.v.pageData.pageAddNote.buttonAdd, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.buttonCancel, 0, 1, 1, 1, 0, 0, true)

	//! add note page grid

	av.v.pageData.pageAddNote.grid.
		SetBorders(false).
		SetRows(0, 1).
		SetColumns(0).
		AddItem(av.v.pageData.pageAddNote.gridData, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageAddNote.gridButtons, 1, 0, 1, 1, 0, 0, true)

	av.v.pageData.pageAddNote.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.pageAddNote.textareaNote:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageAddNote.textareaDesc)
			case av.v.pageData.pageAddNote.textareaDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageAddNote.textareaTitle)
			case av.v.pageData.pageAddNote.textareaTitle:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageAddNote.buttonAdd)
			case av.v.pageData.pageAddNote.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageAddNote.buttonCancel)
			case av.v.pageData.pageAddNote.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageAddNote.textareaNote)
			}
			return nil
		}
		return event
	}

	//! MAIN END -------------------------------------------------------------------

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
				av.v.pageApp.app.SetFocus(av.v.pageData.pageViewNote.textviewNote)
			case av.v.pageData.pageViewNote.textviewNote:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageViewNote.textviewDesc)
			case av.v.pageData.pageViewNote.textviewDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}

	//! adding pages
	av.v.pageData.pages.AddPage("data_add_note", av.v.pageData.pageAddNote.grid, true, true)
	av.v.pageData.pages.AddPage("data_view_note_page", av.v.pageData.pageViewNote.gridData, true, true)
}
