package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PageSelect struct {
	GridSelectButtons *tview.Grid
	ButtonSelect      *tview.Button
	ButtonNew         *tview.Button
	ButtonSettings    *tview.Button
	ButtonDelete      *tview.Button
	ButtonLogout      *tview.Button
	ButtonExit        *tview.Button
	Grid              *tview.Grid
	InputCapture      func(event *tcell.EventKey) *tcell.EventKey
	Page              *tview.Pages
}

type PageEdit struct {
	FormAddEmail      FormAdd
	ButtonAddEmail    *tview.Button
	ButtonDeleteEmail *tview.Button
	ButtonEditExit    *tview.Button
	Grid              *tview.Grid
	InputCapture      func(event *tcell.EventKey) *tcell.EventKey
	Page              *tview.Pages
}

type FormAdd struct {
	InputEmail *tview.InputField
	Form       *tview.Form
}

type PageGroups struct {
	ListGroups *tview.List
	ListEmails *tview.List
	GridMain   *tview.Grid
	PageSelect PageSelect
	PageEdit   PageEdit
	PagesSelEd *tview.Pages
}
