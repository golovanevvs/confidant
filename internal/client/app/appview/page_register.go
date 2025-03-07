package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FormRegisterPage struct {
	Form           *tview.Form
	InputEmail     *tview.InputField
	InputPassword  *tview.InputField
	InputRPassword *tview.InputField
}

type RegisterPage struct {
	InputCapture   func(event *tcell.EventKey) *tcell.EventKey
	Form           FormRegisterPage
	ButtonRegister *tview.Button
	ButtonExit     *tview.Button
	Grid           *tview.Grid
	MainGrid       *tview.Grid
}
