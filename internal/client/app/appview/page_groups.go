package appview

import (
	"context"
	"fmt"

	"github.com/rivo/tview"
)

type pageGroups struct {
	listGroups           *tview.List
	listEmails           *tview.List
	gridMain             *tview.Grid
	pageGroupsSelect     *pageGroupsSelect
	pageGroupsAddGroup   *pageGroupsAddGroup
	pageGroupsEditEmails *pageGroupsEditEmails
	pages                *tview.Pages
}

func newPageGroups() *pageGroups {
	return &pageGroups{
		listGroups:           tview.NewList(),
		listEmails:           tview.NewList(),
		gridMain:             tview.NewGrid(),
		pageGroupsSelect:     newPageGroupsSelect(),
		pageGroupsAddGroup:   newGroupsAddGroup(),
		pageGroupsEditEmails: newPageFroupsEditEmails(),
		pages:                tview.NewPages(),
	}
}

func (av *appView) vGroups() {
	//! Groups List
	av.v.pageGroups.listGroups.ShowSecondaryText(false)
	av.v.pageGroups.listGroups.SetBorder(true)
	av.v.pageGroups.listGroups.SetHighlightFullLine(true)
	av.v.pageGroups.listGroups.SetTitle(" Список групп ")

	av.v.pageGroups.listGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()

		// get data titles list
		var err error
		av.dataTitles, err = av.sv.GetDataTitles(context.Background(), av.account.ID, mainText)
		if err != nil {
			av.v.pageMain.messageBoxL.SetText(err.Error())
		} else {
			av.v.pageData.listTitles.Clear()
			if len(av.dataTitles) > 0 {
				for _, dataTitle := range av.dataTitles {
					av.v.pageData.listTitles.AddItem(dataTitle, "", 0, nil)
				}
			}
			av.titleGroup = mainText
			av.v.pageMain.messageBoxL.SetText(fmt.Sprintf("index: %d, mainText: %s, len: %d, accountID: %d", index, mainText, len(av.dataTitles), av.account.ID))
			av.v.pageData.listTitles.SetTitle(fmt.Sprintf(" %s ", mainText))
			av.v.pageMain.pages.SwitchToPage("data_page")
			av.v.pageData.pages.SwitchToPage("data_view_note_page")
			av.v.pageApp.app.SetInputCapture(av.v.pageData.inputCapture)
			av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
		}
	})

	av.v.pageGroups.listGroups.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// updating e-mails
		av.v.pageGroups.listEmails.Clear()
		for _, email := range av.groups[0].Emails {
			av.v.pageGroups.listEmails.AddItem(email, "", 0, nil)
		}
	})

	//! Main grid
	av.v.pageGroups.gridMain.
		SetRows(0).
		SetColumns(0, 30, 20, 30, 0).
		SetGap(1, 1).
		AddItem(av.v.pageGroups.listGroups, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pages, 0, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.listEmails, 0, 3, 1, 1, 0, 0, true)

	//! adding pages
	av.v.pageGroups.pages.AddPage("select_page", av.v.pageGroups.pageGroupsSelect.grid, true, true)
	av.v.pageGroups.pages.AddPage("add_group_page", av.v.pageGroups.pageGroupsAddGroup.grid, true, true)
	av.v.pageGroups.pages.AddPage("edit_emails_page", av.v.pageGroups.pageGroupsEditEmails.grid, true, true)
}
