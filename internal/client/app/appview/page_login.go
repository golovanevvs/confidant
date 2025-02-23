package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (av *AppView) LoginPage(app *tview.Application, pages *tview.Pages) tview.Primitive {
	// var email, password string
	form := tview.NewForm()
	form.SetHorizontal(false)
	form.AddInputField("E-mail:", "", 0, nil, nil)
	form.AddPasswordField("Пароль:", "", 0, '*', nil)
	// form.AddButton("Войти", func() {
	// 	email = form.GetFormItem(0).(*tview.InputField).GetText()
	// 	password = form.GetFormItem(1).(*tview.InputField).GetText()
	// 	modal := tview.NewModal()
	// 	modal.SetText(fmt.Sprintf("E-mail: %s\nПароль: %s", email, password))
	// 	modal.AddButtons([]string{"Закрыть"})
	// 	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	// 		pages.SwitchToPage("login_page")
	// 	})

	// 	pages.AddPage("modal", modal, true, true)
	// })
	// form.AddButton("Зарегистрироваться", nil)
	// form.AddButton("Выход", func() {
	// 	app.Stop()
	// })

	// form.SetBorder(true).SetTitle("Вход")

	loginButton := tview.NewButton("Войти")
	registerButton := tview.NewButton("Зарегистрироваться").SetSelectedFunc(func() {
		pages.AddAndSwitchToPage("register_page", av.RegisterPage(app, pages), true)
	})
	exitButton := tview.NewButton("Выход").SetSelectedFunc(func() {
		app.Stop()
	})

	formGrid := tview.NewGrid().
		SetRows(5, 3, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(form, 0, 0, 1, 1, 0, 0, true).
		AddItem(loginButton, 1, 0, 1, 1, 1, 1, false).
		AddItem(registerButton, 2, 0, 1, 1, 0, 0, false).
		AddItem(exitButton, 3, 0, 1, 1, 0, 0, false)

	mainGrid := tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(formGrid, 1, 1, 1, 1, 0, 0, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case form.GetFormItem(0):
				app.SetFocus(form.GetFormItem(1))
			case form.GetFormItem(1):
				app.SetFocus(loginButton)
			case loginButton:
				app.SetFocus(registerButton)
			case registerButton:
				app.SetFocus(exitButton)
			case exitButton:
				app.SetFocus(form.GetFormItem(0))
			}
			return nil
		}

		if event.Key() == tcell.KeyEnter {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case form.GetFormItem(1):
				app.SetFocus(loginButton)
			}
		}

		return event
	})

	return mainGrid
}
