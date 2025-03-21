package appview

import (
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
	for i := 0; i < 10; i++ {
		av.v.pageGroups.listGroups.AddItem(fmt.Sprintf("Group %d", i), "", 0, nil)
	}

	av.v.pageGroups.listGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText + secondaryText + string(shortcut))
		av.v.pageData.listTitles.SetTitle(fmt.Sprintf(" %s ", mainText))
		av.v.pageMain.pages.SwitchToPage("data_page")
		av.v.pageData.pages.SwitchToPage("data_view_note_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
	})

	av.v.pageGroups.listGroups.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.v.pageMain.messageBoxL.SetText(mainText)
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
