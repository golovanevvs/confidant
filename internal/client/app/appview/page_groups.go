package appview

import (
	"context"

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

	//! groups list selected
	av.v.pageGroups.listGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		av.groupTitle = mainText
		av.aPageDataSwitch()
	})

	//! group list changed
	av.v.pageGroups.listGroups.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// updating e-mails
		av.aPageGroupsUpdateListEmails(index)
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

func (av *appView) aPageGroupsUpdateListEmails(groupIndex int) {
	av.v.pageGroups.listEmails.Clear()
	for _, email := range av.groups[groupIndex].Emails {
		av.v.pageGroups.listEmails.AddItem(email, "", 0, nil)
	}
}

func (av *appView) aPageGroupsUpdateListGroups() {
	// updating groups list
	av.v.pageGroups.listGroups.Clear()
	if len(av.groups) > 0 {
		for _, group := range av.groups {
			av.v.pageGroups.listGroups.AddItem(group.Title, "", 0, nil)
		}
		// updating e-mails
		av.aPageGroupsUpdateListEmails(0)
	}
}

func (av *appView) aPageGroupsSwitch() {
	var err error
	// getting available groups
	av.groups, err = av.sv.GetGroups(context.Background(), av.account.Email)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(err.Error())
	} else {
		av.aPageGroupsUpdateListGroups()
		av.v.pageMain.pages.SwitchToPage("groups_page")
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
	}
}
