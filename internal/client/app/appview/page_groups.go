package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GroupsPage struct {
	GroupsTable  *tview.Table
	ButtonSelect *tview.Button
	ButtonNew    *tview.Button
	ButtonLogout *tview.Button
	ButtonExit   *tview.Button
	Grid         *tview.Grid
	MainGrid     *tview.Grid
	InputCapture func(event *tcell.EventKey) *tcell.EventKey
}
