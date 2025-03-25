package appview

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageGroupsAddGroup struct {
	formGroupsAddGroup *formGroupsAddGroup
	buttonNew          *tview.Button
	buttonExit         *tview.Button
	grid               *tview.Grid
	inputCapture       func(event *tcell.EventKey) *tcell.EventKey
	page               *tview.Pages
}

type formGroupsAddGroup struct {
	inputName *tview.InputField
	form      *tview.Form
}

func newGroupsAddGroup() *pageGroupsAddGroup {
	return &pageGroupsAddGroup{
		formGroupsAddGroup: &formGroupsAddGroup{
			inputName: tview.NewInputField(),
			form:      tview.NewForm(),
		},
		buttonNew:  tview.NewButton("Создать группу"),
		buttonExit: tview.NewButton("Назад"),
		grid:       tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vGroupsAddGroup() {
	//! add form
	av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.form.SetHorizontal(false)
	av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.form.AddInputField("", "", 0, nil, nil)
	av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.inputName = av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.form.GetFormItem(0).(*tview.InputField)

	//! "Создать группу"
	av.v.pageGroups.pageGroupsAddGroup.buttonNew.SetSelectedFunc(func() {
		title := av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.inputName.GetText()
		av.sv.AddGroup(context.Background(), &av.account, title)
	})

	//! "Назад"
	av.v.pageGroups.pageGroupsAddGroup.buttonExit.SetSelectedFunc(func() {
		av.v.pageGroups.pages.SwitchToPage("select_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsSelect.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
	})

	//! "buttons grid"
	av.v.pageGroups.pageGroupsAddGroup.grid.
		SetRows(3, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsAddGroup.buttonNew, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageGroups.pageGroupsAddGroup.buttonExit, 2, 0, 1, 1, 0, 0, true)

	//! InputCapture add group page
	av.v.pageGroups.pageGroupsAddGroup.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.inputName:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsAddGroup.buttonNew)
			case av.v.pageGroups.pageGroupsAddGroup.buttonNew:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsAddGroup.buttonExit)
			case av.v.pageGroups.pageGroupsAddGroup.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageGroups.pageGroupsAddGroup.formGroupsAddGroup.inputName)
			}
			return nil
		}
		return event
	}
}
