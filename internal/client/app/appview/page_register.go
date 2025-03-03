package appview

// func (av *AppView) RegisterPage(app *tview.Application, pages *tview.Pages) tview.Primitive {
// var email, password string
// form := tview.NewForm()
// form.SetHorizontal(false)
// form.AddInputField("E-mail:", "", 0, nil, nil)
// form.AddPasswordField("Пароль:", "", 0, '*', nil)
// form.AddPasswordField("Повторите:", "", 0, '*', nil)

// registerButton := tview.NewButton("Зарегистрироваться").
// 	SetSelectedFunc(func() {
// 		//! Зарегистрироваться
// 		pass1 := form.GetFormItem(1).(*tview.InputField).GetText()
// 		pass2 := form.GetFormItem(2).(*tview.InputField).GetText()
// 		if pass1 != pass2 {
// 			modal := tview.NewModal().
// 				SetText("Пароли не совпадают! Повторите ввод").
// 				AddButtons([]string{"Закрыть"}).
// 				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
// 					if buttonLabel == "Закрыть" {
// 						pages.SwitchToPage("register_page")
// 					}
// 				})
// 			pages.AddPage("register_page_modal", modal, true, true)
// 		} else {
// 			account := model.Account{
// 				Email:    form.GetFormItem(0).(*tview.InputField).GetText(),
// 				Password: form.GetFormItem(1).(*tview.InputField).GetText(),
// 			}
// 			tr := trhttp.New()
// 			tr.RegisterAccount(account)
// 		}
// 	})
// exitButton := tview.NewButton("Назад").SetSelectedFunc(func() {
// 	pages.AddAndSwitchToPage("login_page", av.LoginPage(app, pages), true)
// })

// formGrid := tview.NewGrid().
// 	SetRows(8, 3, 3).
// 	SetColumns(50).
// 	SetGap(1, 0).
// 	AddItem(form, 0, 0, 1, 1, 0, 0, true).
// 	AddItem(registerButton, 1, 0, 1, 1, 0, 0, true).
// 	AddItem(exitButton, 2, 0, 1, 1, 0, 0, true)

// mainGrid := tview.NewGrid().
// 	SetRows(0, 20, 0).
// 	SetColumns(0, 40, 0).
// 	AddItem(formGrid, 1, 1, 1, 1, 0, 0, true)

// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 	if event.Key() == tcell.KeyTAB {
// 		currentFocus := app.GetFocus()
// 		switch currentFocus {
// 		case form.GetFormItem(0):
// 			app.SetFocus(form.GetFormItem(1))
// 		case form.GetFormItem(1):
// 			app.SetFocus(form.GetFormItem(2))
// 		case form.GetFormItem(2):
// 			app.SetFocus(registerButton)
// 		case registerButton:
// 			app.SetFocus(exitButton)
// 		case exitButton:
// 			app.SetFocus(form.GetFormItem(0))
// 		}
// 		return nil
// 	}

// 	if event.Key() == tcell.KeyEnter {
// 		currentFocus := app.GetFocus()
// 		switch currentFocus {
// 		case form.GetFormItem(2):
// 			app.SetFocus(registerButton)
// 		}
// 	}

// 	return event
// })

// 	return mainGrid
// }
