package appview

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageData struct {
	listTitles         *tview.List
	buttonAdd          *tview.Button
	buttonEdit         *tview.Button
	buttonDelete       *tview.Button
	buttonBack         *tview.Button
	buttonExit         *tview.Button
	gridButtons        *tview.Grid
	gridMain           *tview.Grid
	pageDataSelectType *pageDataSelectType
	pageDataViewEmpty  *pageDataViewEmpty
	pageDataViewNote   *pageDataViewNote
	pageDataViewPass   *pageDataViewPass
	pageDataViewCard   *pageDataViewCard
	pageDataViewFile   *pageDataViewFile
	pageDataAddNote    *pageDataAddNote
	pageDataAddPass    *pageDataAddPass
	pageDataAddCard    *pageDataAddCard
	pageDataAddFile    *pageDataAddFile
	pages              *tview.Pages
	inputCapture       func(event *tcell.EventKey) *tcell.EventKey
	page               *tview.Pages
}

func newPageData() *pageData {
	return &pageData{
		listTitles:         tview.NewList(),
		buttonAdd:          tview.NewButton("Добавить"),
		buttonEdit:         tview.NewButton("Изменить"),
		buttonDelete:       tview.NewButton("Удалить"),
		buttonBack:         tview.NewButton("Назад"),
		buttonExit:         tview.NewButton("Выход"),
		gridButtons:        tview.NewGrid(),
		gridMain:           tview.NewGrid(),
		pageDataSelectType: newPageDataSelectType(),
		pageDataViewEmpty:  newPageDataViewEmpty(),
		pageDataViewNote:   newPageDataViewNote(),
		pageDataViewPass:   newPageDataViewPass(),
		pageDataViewCard:   newPageDataViewCard(),
		pageDataViewFile:   newDataViewFile(),
		pageDataAddNote:    newPageDataAddNote(),
		pageDataAddPass:    newPageDataAddPass(),
		pageDataAddCard:    newPageDataAddCard(),
		pageDataAddFile:    newDataAddFile(),
		pages:              tview.NewPages(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vData() {
	//! titles list
	av.v.pageData.listTitles.ShowSecondaryText(false)
	av.v.pageData.listTitles.SetBorder(true)
	av.v.pageData.listTitles.SetHighlightFullLine(true)
	// for i := 0; i < 10; i++ {
	// 	av.v.pageData.listTitles.AddItem(fmt.Sprintf("%c Title %d", '\U0001F449', i), "", 0, nil)
	// }

	av.v.pageData.listTitles.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
	})

	av.v.pageData.listTitles.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.dataTitle = mainText[strings.Index(mainText, " ")+1:]
		av.aPageDataUpdateDataView()
	})

	//! "Добавить"
	av.v.pageData.buttonAdd.SetSelectedFunc(func() {
		av.v.pageData.pages.SwitchToPage("data_select_type")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataSelectType.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.pageDataSelectType.buttonNote)
	})

	//! "Назад"
	av.v.pageData.buttonBack.SetSelectedFunc(func() {
		av.vClearMessages()
		av.groupID = -1
		av.aPageGroupsSwitch()
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

	//! adding pages
	av.v.pageData.pages.AddPage("data_select_type", av.v.pageData.pageDataSelectType.grid, true, true)
	av.v.pageData.pages.AddPage("data_view_empty_page", av.v.pageData.pageDataViewEmpty.gridData, true, true)
	av.v.pageData.pages.AddPage("data_view_note_page", av.v.pageData.pageDataViewNote.gridData, true, true)
	av.v.pageData.pages.AddPage("data_view_pass_page", av.v.pageData.pageDataViewPass.gridData, true, true)
	av.v.pageData.pages.AddPage("data_view_card_page", av.v.pageData.pageDataViewCard.gridData, true, true)
	av.v.pageData.pages.AddPage("data_view_file_page", av.v.pageData.pageDataViewFile.grid, true, true)
	av.v.pageData.pages.AddPage("data_add_note_page", av.v.pageData.pageDataAddNote.grid, true, true)
	av.v.pageData.pages.AddPage("data_add_pass_page", av.v.pageData.pageDataAddPass.grid, true, true)
	av.v.pageData.pages.AddPage("data_add_card_page", av.v.pageData.pageDataAddCard.grid, true, true)

	av.v.pageData.pages.AddPage("data_add_file_page", av.v.pageData.pageDataAddFile.grid, true, true)
}

func (av *appView) aPageDataSwitch() {
	// clearing messages
	av.vClearMessages()

	av.v.pageData.listTitles.SetTitle(fmt.Sprintf(" %s ", av.groupTitle))

	var err error

	// getting group ID
	av.groupID, err = av.sv.GetGroupID(context.Background(), av.account.Email, av.groupTitle)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(err.Error())
	} else {
		av.v.pageMain.pages.SwitchToPage("data_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)

		av.aPageDataUpdateListTitles()
	}
}

func (av *appView) aPageDataUpdateListTitles() {
	var err error

	// getting data titles
	av.dataTitles, err = av.sv.GetDataTitles(context.Background(), av.account.ID, av.groupID)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(err.Error())
	} else {
		av.v.pageData.listTitles.Clear()
		av.dataTypes, err = av.sv.GetDataTypes(context.Background(), av.account.ID, av.groupID)
		if err != nil {
			av.v.pageMain.messageBoxL.SetText(err.Error())
		} else {
			av.v.pageData.listTitles.Clear()

			if len(av.dataTitles) > 0 {
				//filling data list
				for i, dataTitle := range av.dataTitles {
					switch av.dataTypes[i] {
					case "note":
						dataTitle = fmt.Sprintf("📝 %s", dataTitle)
					case "pass":
						dataTitle = fmt.Sprintf("🔒 %s", dataTitle)
					case "card":
						dataTitle = fmt.Sprintf("💳 %s", dataTitle)
					case "file":
						dataTitle = fmt.Sprintf("📁 %s", dataTitle)
					}
					av.v.pageData.listTitles.AddItem(dataTitle, "", 0, nil)

				}
				av.dataTitle = av.dataTitles[0]
				av.aPageDataUpdateDataView()
			} else {
				av.aPageDataViewEmptySwitch()
			}

		}
	}
}

func (av *appView) aPageDataUpdateDataView() {
	var err error

	// getting data ID and data type
	av.dataID, av.dataType, err = av.sv.GetDataIDAndType(context.Background(), av.groupID, av.dataTitle)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(fmt.Sprintf("groupID: %d, dataTitle: %s, err: %s", av.groupID, av.dataTitle, err.Error()))
	} else {
		switch av.dataType {
		case "note":
			av.vPageDataViewNoteUpdate()
		case "pass":
			av.vPageDataViewPassUpdate()
		case "card":
			av.vPageDataViewCardUpdate()
		}
	}
}
