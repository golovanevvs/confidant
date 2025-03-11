package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormPageRegister struct {
	Form           *tview.Form
	InputEmail     *tview.InputField
	InputPassword  *tview.InputField
	InputRPassword *tview.InputField
}

type PageRegister struct {
	InputCapture   func(event *tcell.EventKey) *tcell.EventKey
	Form           FormPageRegister
	ButtonRegister *tview.Button
	ButtonExit     *tview.Button
	Grid           *tview.Grid
	MainGrid       *tview.Grid
}
