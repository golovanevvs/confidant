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

func (av *AppView) VLogin() {

	av.v.pageLogin.Form.Form.SetHorizontal(false)
	av.v.pageLogin.Form.Form.AddInputField("E-mail:", "", 0, nil, nil)
	av.v.pageLogin.Form.InputEmail = av.v.pageLogin.Form.Form.GetFormItem(0).(*tview.InputField)
	av.v.pageLogin.Form.Form.AddPasswordField("Пароль:", "", 0, '*', nil)
	av.v.pageLogin.Form.InputPassword = av.v.pageLogin.Form.Form.GetFormItem(1).(*tview.InputField)

	//? Войти
	av.v.pageLogin.ButtonLogin.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.Pages.SwitchToPage("groups_page")
		av.v.pageGroups.PagesSelEd.SwitchToPage("select_page")
		av.v.pageMain.App.SetInputCapture(av.v.pageGroups.PageSelect.InputCapture)
		av.v.pageMain.App.SetFocus(av.v.pageGroups.ListGroups)
	})

	//? Регистрация
	av.v.pageLogin.ButtonRegister.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.Pages.SwitchToPage("register_page")
		// focus
		av.v.pageMain.App.SetInputCapture(av.v.pageRegister.InputCapture)
		av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputEmail)
		// clear
		av.v.pageRegister.Form.InputEmail.SetText("")
		av.v.pageRegister.Form.InputPassword.SetText("")
		av.v.pageRegister.Form.InputRPassword.SetText("")
		// messageBox
		av.v.pageMain.MessageBoxL.SetText("Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
	})

	//? Выход
	av.v.pageLogin.ButtonExit.SetSelectedFunc(func() {
		av.v.pageMain.App.Stop()
	})

	//? form grid
	av.v.pageLogin.Grid.
		SetRows(5, 1, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(av.v.pageLogin.Form.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageLogin.ButtonLogin, 1, 0, 1, 1, 1, 1, false).
		AddItem(av.v.pageLogin.ButtonRegister, 2, 0, 1, 1, 0, 0, false).
		AddItem(av.v.pageLogin.ButtonExit, 3, 0, 1, 1, 0, 0, false)

	//? main grid
	av.v.pageLogin.MainGrid.
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(av.v.pageLogin.Grid, 1, 1, 1, 1, 0, 0, true)

	//? InputCapture
	av.v.pageLogin.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageMain.App.GetFocus()
			switch currentFocus {
			case av.v.pageLogin.Form.InputEmail:
				av.v.pageMain.App.SetFocus(av.v.pageLogin.Form.InputPassword)
			case av.v.pageLogin.Form.InputPassword:
				av.v.pageMain.App.SetFocus(av.v.pageLogin.ButtonLogin)
			case av.v.pageLogin.ButtonLogin:
				av.v.pageMain.App.SetFocus(av.v.pageLogin.ButtonRegister)
			case av.v.pageLogin.ButtonRegister:
				av.v.pageMain.App.SetFocus(av.v.pageLogin.ButtonExit)
			case av.v.pageLogin.ButtonExit:
				av.v.pageMain.App.SetFocus(av.v.pageLogin.Form.InputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := av.v.pageMain.App.GetFocus()
		// 	switch currentFocus {
		// 	case formav.v.pageLogin.GetFormItem(1):
		// 		av.v.pageMain.App.SetFocus(buttonLoginav.v.pageLogin)
		// 	}
		// }

		return event
	}
}
