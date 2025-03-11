package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormPageLogin struct {
	Form          *tview.Form
	InputEmail    *tview.InputField
	InputPassword *tview.InputField
}

type PageLogin struct {
	InputCapture   func(event *tcell.EventKey) *tcell.EventKey
	Form           FormPageLogin
	ButtonLogin    *tview.Button
	ButtonRegister *tview.Button
	ButtonExit     *tview.Button
	Grid           *tview.Grid
	MainGrid       *tview.Grid
}
