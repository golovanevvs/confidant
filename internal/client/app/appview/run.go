package appview

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

func (av *AppView) Run() error {
	//!Beginning
	action := "new"
	// app
	app := tview.NewApplication()
	app.EnableMouse(true)

	// page container
	pages := tview.NewPages()

	// left message box
	messageBoxL := tview.NewTextView()
	messageBoxL.SetDynamicColors(true)
	messageBoxL.SetTextAlign(tview.AlignLeft)
	messageBoxL.SetText("[green]Добро пожаловать в систему хранения конфиденциальной информации [white]CON[blue]FID[red]ANT")
	messageBoxL.SetBorder(true)
	messageBoxL.SetBorderColor(tcell.ColorRed)
	messageBoxL.SetTitle(" Сообщения ")

	// right message box
	messageBoxR := tview.NewTextView()
	messageBoxR.SetDynamicColors(true)
	messageBoxR.SetTextAlign(tview.AlignLeft)
	messageBoxR.SetBorder(true)
	messageBoxR.SetBorderColor(tcell.ColorRed)
	messageBoxR.SetTitle(" Дополнительная информация ")

	// main grid
	mainGrid := tview.NewGrid()
	mainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	mainGrid.SetRows(0, 8)
	mainGrid.AddItem(pages, 0, 0, 1, 2, 0, 0, true)
	mainGrid.AddItem(messageBoxL, 1, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(messageBoxR, 1, 1, 1, 1, 0, 0, true)

	loginPage := LoginPage{}
	registerPage := RegisterPage{}
	groupsPage := GroupsPage{}

	//! LOGIN PAGE

	loginPage.Form.Form = tview.NewForm()
	loginPage.Form.Form.SetHorizontal(false)
	loginPage.Form.Form.AddInputField("E-mail:", "", 0, nil, nil)
	loginPage.Form.InputEmail = loginPage.Form.Form.GetFormItem(0).(*tview.InputField)
	loginPage.Form.Form.AddPasswordField("Пароль:", "", 0, '*', nil)
	loginPage.Form.InputPassword = loginPage.Form.Form.GetFormItem(1).(*tview.InputField)
	//! Войти
	loginPage.ButtonLogin = tview.NewButton("Войти").SetSelectedFunc(func() {
		// switch
		pages.SwitchToPage("groups_page")
		app.SetInputCapture(groupsPage.InputCapture)
	})
	//! Зарегистрироваться
	loginPage.ButtonRegister = tview.NewButton("Зарегистрироваться").SetSelectedFunc(func() {
		// switch
		pages.SwitchToPage("register_page")
		// focus
		app.SetInputCapture(registerPage.InputCapture)
		app.SetFocus(registerPage.Form.InputEmail)
		// clear
		registerPage.Form.InputEmail.SetText("")
		registerPage.Form.InputPassword.SetText("")
		registerPage.Form.InputRPassword.SetText("")
		// messageBox
		messageBoxL.SetText("Введите e-mail и пароль. Повторите ввод пароля. Нажмите кнопку [blue]\"Зарегистрироваться\".\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
	})
	//! Выход
	loginPage.ButtonExit = tview.NewButton("Выход").SetSelectedFunc(func() {
		app.Stop()
	})

	loginPage.Grid = tview.NewGrid().
		SetRows(5, 3, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(loginPage.Form.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(loginPage.ButtonLogin, 1, 0, 1, 1, 1, 1, false).
		AddItem(loginPage.ButtonRegister, 2, 0, 1, 1, 0, 0, false).
		AddItem(loginPage.ButtonExit, 3, 0, 1, 1, 0, 0, false)

	loginPage.MainGrid = tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(loginPage.Grid, 1, 1, 1, 1, 0, 0, true)

	loginPage.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case loginPage.Form.InputEmail:
				app.SetFocus(loginPage.Form.InputPassword)
			case loginPage.Form.InputPassword:
				app.SetFocus(loginPage.ButtonLogin)
			case loginPage.ButtonLogin:
				app.SetFocus(loginPage.ButtonRegister)
			case loginPage.ButtonRegister:
				app.SetFocus(loginPage.ButtonExit)
			case loginPage.ButtonExit:
				app.SetFocus(loginPage.Form.InputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := app.GetFocus()
		// 	switch currentFocus {
		// 	case formLoginPage.GetFormItem(1):
		// 		app.SetFocus(buttonLoginLoginPage)
		// 	}
		// }

		return event
	}

	//! REGISTER PAGE

	//! Form
	registerPage.Form.Form = tview.NewForm()
	registerPage.Form.Form.SetHorizontal(false)
	registerPage.Form.Form.AddInputField("E-mail:", "", 0, nil, nil)
	registerPage.Form.InputEmail = registerPage.Form.Form.GetFormItem(0).(*tview.InputField)
	registerPage.Form.Form.AddPasswordField("Пароль:", "", 0, '*', nil)
	registerPage.Form.InputPassword = registerPage.Form.Form.GetFormItem(1).(*tview.InputField)
	registerPage.Form.Form.AddPasswordField("Повторите:", "", 0, '*', nil)
	registerPage.Form.InputRPassword = registerPage.Form.Form.GetFormItem(2).(*tview.InputField)

	registerPage.ButtonRegister = tview.NewButton("Зарегистрироваться").
		SetSelectedFunc(func() {
			//! Зарегистрироваться
			pass1 := registerPage.Form.InputPassword.GetText()
			pass2 := registerPage.Form.InputRPassword.GetText()
			if pass1 == pass2 {
				email := registerPage.Form.InputEmail.GetText()
				password := registerPage.Form.InputPassword.GetText()
				registerAccountResp, err := av.sv.RegisterAccount(email, password)
				if err != nil {
					messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
					messageBoxL.SetText("[red]Возникла критическая ошибка.")
				} else {
					if registerAccountResp.HTTPStatusCode != 200 {
						messageBoxR.SetText(fmt.Sprintf("[red]%s", registerAccountResp.ServerError))
						switch {
						// invalid e-mail
						case strings.Contains(registerAccountResp.ServerError, customerrors.ErrAccountValidateEmail422.Error()):
							messageBoxL.SetText("[red]Неверно введён e-mail!")
							app.SetFocus(registerPage.Form.InputEmail)
						// invalid password
						case strings.Contains(registerAccountResp.ServerError, customerrors.ErrAccountValidatePassword422.Error()):
							messageBoxL.SetText("[red]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов!")
							registerPage.Form.InputPassword.SetText("")
							registerPage.Form.InputRPassword.SetText("")
							app.SetFocus(registerPage.Form.InputPassword)
						// e-mail is already busy
						case strings.Contains(registerAccountResp.ServerError, customerrors.ErrDBBusyEmail409.Error()):
							messageBoxL.SetText(fmt.Sprintf("[red]Пользователь с e-mail %s уже зарегестрирован!", email))
							app.SetFocus(registerPage.Form.InputEmail)
						// other errors
						default:
							messageBoxL.SetText("[red]Возникла ошибка.")
						}
					} else {
						messageBoxR.Clear()
						messageBoxL.SetText(fmt.Sprintf("[green]Вы успешно зарегистрировались. Ваш ID: %s\n[white]Войдите в систему, используя свой e-mail и пароль.", registerAccountResp.AccountID))
					}
				}
			} else {
				messageBoxL.Clear()
				messageBoxL.SetText("[red]Пароли не совпадают! Повторите ввод.\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
				messageBoxR.Clear()
				registerPage.Form.InputPassword.SetText("")
				registerPage.Form.InputRPassword.SetText("")
				app.SetFocus(registerPage.Form.InputPassword)
			}
		})
	//! Назад
	registerPage.ButtonExit = tview.NewButton("Назад").SetSelectedFunc(func() {
		// switch
		// pages.AddAndSwitchToPage("login_page", mainGridLoginPage, true)
		pages.SwitchToPage("login_page")
		// focus
		app.SetInputCapture(loginPage.InputCapture)
		app.SetFocus(loginPage.Form.InputEmail)
		// messageBox
		messageBoxL.SetText("[green]Добро пожаловать в систему хранения конфиденциальной информации [white]CON[blue]FID[red]ANT")
		messageBoxR.Clear()
	})

	//! Grid
	registerPage.Grid = tview.NewGrid().
		SetRows(8, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(registerPage.Form.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(registerPage.ButtonRegister, 1, 0, 1, 1, 0, 0, true).
		AddItem(registerPage.ButtonExit, 2, 0, 1, 1, 0, 0, true)

	//! MainGrid
	registerPage.MainGrid = tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(registerPage.Grid, 1, 1, 1, 1, 0, 0, true)

	//! InputCapture
	registerPage.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case registerPage.Form.InputEmail:
				app.SetFocus(registerPage.Form.InputPassword)
			case registerPage.Form.InputPassword:
				app.SetFocus(registerPage.Form.InputRPassword)
			case registerPage.Form.InputRPassword:
				app.SetFocus(registerPage.ButtonRegister)
			case registerPage.ButtonRegister:
				app.SetFocus(registerPage.ButtonExit)
			case registerPage.ButtonExit:
				app.SetFocus(registerPage.Form.InputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := app.GetFocus()
		// 	switch currentFocus {
		// 	case registerPage.Form.Form.GetFormItem(2):
		// 		app.SetFocus(registerPage.ButtonRegister)
		// 	}
		// }

		return event
	}

	//! GROUPS PAGE

	//! Groups List
	groupsPage.GroupsList = tview.NewList()
	groupsPage.GroupsList.SetBorder(true)
	for i := 0; i < 10; i++ {
		groupsPage.GroupsList.AddItem(fmt.Sprintf("Group %d", i), "", 0, nil)
	}

	groupsPage.GroupsList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		messageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	//! Выбрать группу
	groupsPage.ButtonSelect = tview.NewButton("Выбрать группу")

	//! Создать группу
	groupsPage.ButtonNew = tview.NewButton("Создать группу")

	//! Удалить группу
	groupsPage.ButtonDelete = tview.NewButton("Удалить группу")

	//! Выйти из аккаунта
	groupsPage.ButtonLogout = tview.NewButton("Выйти из аккаунта")

	//! Выход
	groupsPage.ButtonExit = tview.NewButton("Выход")

	//! Grid
	groupsPage.Grid = tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(groupsPage.ButtonSelect, 1, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.ButtonNew, 2, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.ButtonDelete, 3, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.ButtonLogout, 4, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.ButtonExit, 5, 0, 1, 1, 0, 0, true)

	//! Main grid
	groupsPage.MainGrid = tview.NewGrid().
		SetRows(0).
		SetColumns(40, 40, 0).
		SetGap(1, 1).
		AddItem(groupsPage.GroupsList, 0, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.Grid, 0, 1, 1, 1, 0, 0, true)

	//! InputCapture
	groupsPage.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case groupsPage.GroupsList:
				app.SetFocus(groupsPage.ButtonSelect)
			case groupsPage.ButtonSelect:
				app.SetFocus(groupsPage.ButtonNew)
			case groupsPage.ButtonNew:
				app.SetFocus(groupsPage.ButtonDelete)
			case groupsPage.ButtonDelete:
				app.SetFocus(groupsPage.ButtonLogout)
			case groupsPage.ButtonLogout:
				app.SetFocus(groupsPage.ButtonExit)
			case groupsPage.ButtonExit:
				app.SetFocus(groupsPage.GroupsList)
			}
			return nil
		}
		return event
	}

	//! Adding Pages

	pages.AddPage("groups_page", groupsPage.MainGrid, true, true)
	pages.AddPage("register_page", registerPage.MainGrid, true, true)
	pages.AddAndSwitchToPage("login_page", loginPage.MainGrid, true)
	app.SetInputCapture(loginPage.InputCapture)
	app.SetFocus(loginPage.Form.InputEmail)

	//! Launching the app
	app.SetRoot(mainGrid, true)

	if err := app.Run(); err != nil {
		return fmt.Errorf("%s: %s: %w: %w", customerrors.ClientAppViewErr, action, customerrors.ErrRunAppView, err)
	}

	return nil
}
