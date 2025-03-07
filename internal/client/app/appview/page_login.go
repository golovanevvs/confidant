package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormLoginPage struct {
	Form          *tview.Form
	InputEmail    *tview.InputField
	InputPassword *tview.InputField
}

type LoginPage struct {
	InputCapture   func(event *tcell.EventKey) *tcell.EventKey
	Form           FormLoginPage
	ButtonLogin    *tview.Button
	ButtonRegister *tview.Button
	ButtonExit     *tview.Button
	Grid           *tview.Grid
	MainGrid       *tview.Grid
}
