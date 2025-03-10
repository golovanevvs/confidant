package appview

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

func (av *AppView) Run() error {
	//!Beginning
	action := "new"

	var statusServer *model.StatusResp
	var statusServerErr error

	// app
	app := tview.NewApplication()
	// app.EnableMouse(true)

	// page container
	pages := tview.NewPages()

	//? left message box
	messageBoxL := tview.NewTextView()
	messageBoxL.SetDynamicColors(true)
	messageBoxL.SetTextAlign(tview.AlignLeft)
	messageBoxL.SetBorder(true).SetBorderColor(tcell.ColorRed)
	messageBoxL.SetTitle(" Сообщения ")

	//? right message box
	messageBoxR := tview.NewTextView()
	messageBoxR.SetDynamicColors(true)
	messageBoxR.SetTextAlign(tview.AlignLeft)
	messageBoxR.SetBorder(true).SetBorderColor(tcell.ColorRed)
	messageBoxR.SetTitle(" Дополнительная информация ")

	//? status bar
	statusBar := tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorRed)
	statusBar.SetCell(0, 0, tview.NewTableCell("Тип соединения").SetAlign(tview.AlignCenter).SetExpansion(1))
	statusBar.SetCell(0, 1, tview.NewTableCell("Соединение с сервером").SetAlign(tview.AlignCenter).SetExpansion(1))
	statusBar.SetCell(0, 2, tview.NewTableCell("Соединение с БД сервера").SetAlign(tview.AlignCenter).SetExpansion(1))
	statusBar.SetCell(0, 3, tview.NewTableCell("Статус операции").SetAlign(tview.AlignCenter).SetExpansion(1))
	statusBar.SetCell(0, 4, tview.NewTableCell("Активный аккаунт").SetAlign(tview.AlignCenter).SetExpansion(1))
	statusBarCellTypeConnect := tview.NewTableCell("[green]REST API").SetAlign(tview.AlignCenter).SetExpansion(1)
	statusBar.SetCell(1, 0, statusBarCellTypeConnect)
	statusBarCellServerConnect := tview.NewTableCell("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	statusBar.SetCell(1, 1, statusBarCellServerConnect)
	statusBarCellServerDBConnect := tview.NewTableCell("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	statusBar.SetCell(1, 2, statusBarCellServerDBConnect)

	updateCellServerConnectChan := make(chan string)
	updateCellServerDBConnectChan := make(chan string)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			statusServer, statusServerErr = av.sv.GetServerStatus()
			if statusServerErr != nil {
				updateCellServerConnectChan <- "[red]Отсутствует"
				updateCellServerDBConnectChan <- "[red]Отсутствует"
			} else {
				if statusServer.HTTPStatusCode == 200 {
					updateCellServerConnectChan <- "[green]OK"
					updateCellServerDBConnectChan <- "[green]OK"
				} else {
					updateCellServerConnectChan <- "[green]OK"
					updateCellServerDBConnectChan <- "[red]Отсутствует"
				}
			}
		}
	}()

	go func() {
		for info := range updateCellServerConnectChan {
			app.QueueUpdateDraw(func() {
				statusBarCellServerConnect.SetText(info)
			})
		}
	}()
	go func() {
		for info := range updateCellServerDBConnectChan {
			app.QueueUpdateDraw(func() {
				statusBarCellServerDBConnect.SetText(info)
			})
		}
	}()

	statusBarCellResponseStatus := tview.NewTableCell("").SetAlign(tview.AlignCenter).SetExpansion(1)
	statusBar.SetCell(1, 3, statusBarCellResponseStatus)
	statusBarCellActiveAccount := tview.NewTableCell("").SetAlign(tview.AlignCenter).SetExpansion(1)
	statusBar.SetCell(1, 4, statusBarCellActiveAccount)

	//? main grid
	mainGrid := tview.NewGrid()
	mainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	mainGrid.SetRows(0, 8, 5)
	mainGrid.AddItem(pages, 0, 0, 1, 2, 0, 0, true)
	mainGrid.AddItem(messageBoxL, 1, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(messageBoxR, 1, 1, 1, 1, 0, 0, true)
	mainGrid.AddItem(statusBar, 2, 0, 1, 2, 0, 0, false)

	//? pages
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

	//? Войти
	loginPage.ButtonLogin = tview.NewButton("Войти").SetSelectedFunc(func() {
		// switch
		pages.SwitchToPage("groups_page")
		groupsPage.PagesSelEd.SwitchToPage("select_page")
		app.SetInputCapture(groupsPage.PageSelect.InputCapture)
		app.SetFocus(groupsPage.ListGroups)
	})

	//? Зарегистрироваться
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

	//? Выход
	loginPage.ButtonExit = tview.NewButton("Выход").SetSelectedFunc(func() {
		app.Stop()
	})

	//? form grid
	loginPage.Grid = tview.NewGrid().
		SetRows(5, 1, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(loginPage.Form.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(loginPage.ButtonLogin, 1, 0, 1, 1, 1, 1, false).
		AddItem(loginPage.ButtonRegister, 2, 0, 1, 1, 0, 0, false).
		AddItem(loginPage.ButtonExit, 3, 0, 1, 1, 0, 0, false)

	//? main grid
	loginPage.MainGrid = tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(loginPage.Grid, 1, 1, 1, 1, 0, 0, true)

	//? InputCapture
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

	//? Form
	registerPage.Form.Form = tview.NewForm()
	registerPage.Form.Form.SetHorizontal(false)
	registerPage.Form.Form.AddInputField("E-mail:", "", 0, nil, nil)
	registerPage.Form.InputEmail = registerPage.Form.Form.GetFormItem(0).(*tview.InputField)
	registerPage.Form.Form.AddPasswordField("Пароль:", "", 0, '*', nil)
	registerPage.Form.InputPassword = registerPage.Form.Form.GetFormItem(1).(*tview.InputField)
	registerPage.Form.Form.AddPasswordField("Повторите:", "", 0, '*', nil)
	registerPage.Form.InputRPassword = registerPage.Form.Form.GetFormItem(2).(*tview.InputField)

	//? Зарегистрироваться
	registerPage.ButtonRegister = tview.NewButton("Зарегистрироваться").
		SetSelectedFunc(func() {
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

	//? Назад
	registerPage.ButtonExit = tview.NewButton("Назад").SetSelectedFunc(func() {
		// switch
		pages.SwitchToPage("login_page")
		// focus
		app.SetInputCapture(loginPage.InputCapture)
		app.SetFocus(loginPage.Form.InputEmail)
		// messageBox
		messageBoxL.Clear()
		messageBoxR.Clear()
	})

	//? form grid
	registerPage.Grid = tview.NewGrid().
		SetRows(8, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(registerPage.Form.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(registerPage.ButtonRegister, 1, 0, 1, 1, 0, 0, true).
		AddItem(registerPage.ButtonExit, 2, 0, 1, 1, 0, 0, true)

	//? main grid
	registerPage.MainGrid = tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(registerPage.Grid, 1, 1, 1, 1, 0, 0, true)

	//? InputCapture
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

	//? pages
	groupsPage.PagesSelEd = tview.NewPages()

	//? Groups List
	groupsPage.ListGroups = tview.NewList()
	groupsPage.ListGroups.SetBorder(true)
	groupsPage.ListGroups.SetHighlightFullLine(true)
	groupsPage.ListGroups.SetTitle(" Список групп ")
	for i := 0; i < 10; i++ {
		groupsPage.ListGroups.AddItem(fmt.Sprintf("Group %d", i), "", 0, nil)
	}

	groupsPage.ListGroups.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		messageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	//? "select_page"

	//? Выбрать группу
	groupsPage.PageSelect.ButtonSelect = tview.NewButton("Выбрать группу")

	//? Создать группу
	groupsPage.PageSelect.ButtonNew = tview.NewButton("Создать группу")

	//? Настроить группу
	groupsPage.PageSelect.ButtonSettings = tview.NewButton("Настроить группу").SetSelectedFunc(func() {
		groupsPage.PagesSelEd.SwitchToPage("edit_page")
		app.SetInputCapture(groupsPage.PageEdit.InputCapture)
		app.SetFocus(groupsPage.ListEmails)
		groupsPage.PageEdit.FormAddEmail.InputEmail.SetText("")
	})

	//? Удалить группу
	groupsPage.PageSelect.ButtonDelete = tview.NewButton("Удалить группу")

	//? Выйти из аккаунта
	groupsPage.PageSelect.ButtonLogout = tview.NewButton("Выйти из аккаунта").SetSelectedFunc(func() {
		// switch
		pages.SwitchToPage("login_page")
		// focus
		app.SetInputCapture(loginPage.InputCapture)
		app.SetFocus(loginPage.Form.InputEmail)
		// messageBox
		messageBoxL.Clear()
		messageBoxR.Clear()
	})

	//? Выход
	groupsPage.PageSelect.ButtonExit = tview.NewButton("Выход").SetSelectedFunc(func() {
		app.Stop()
	})

	//? MainButtonsGrid
	groupsPage.PageSelect.Grid = tview.NewGrid().
		SetRows(1, 1, 1, 1, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(groupsPage.PageSelect.ButtonSelect, 1, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageSelect.ButtonNew, 2, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageSelect.ButtonSettings, 3, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageSelect.ButtonDelete, 4, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageSelect.ButtonLogout, 5, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageSelect.ButtonExit, 6, 0, 1, 1, 0, 0, true)

	//? InputCapture select page
	groupsPage.PageSelect.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case groupsPage.ListGroups:
				app.SetFocus(groupsPage.PageSelect.ButtonSelect)
			case groupsPage.PageSelect.ButtonSelect:
				app.SetFocus(groupsPage.PageSelect.ButtonNew)
			case groupsPage.PageSelect.ButtonNew:
				app.SetFocus(groupsPage.PageSelect.ButtonSettings)
			case groupsPage.PageSelect.ButtonSettings:
				app.SetFocus(groupsPage.PageSelect.ButtonDelete)
			case groupsPage.PageSelect.ButtonDelete:
				app.SetFocus(groupsPage.PageSelect.ButtonLogout)
			case groupsPage.PageSelect.ButtonLogout:
				app.SetFocus(groupsPage.PageSelect.ButtonExit)
			case groupsPage.PageSelect.ButtonExit:
				app.SetFocus(groupsPage.ListGroups)
			}
			return nil
		}
		return event
	}

	//? "edit_page"

	//? add form
	groupsPage.PageEdit.FormAddEmail.Form = tview.NewForm()
	groupsPage.PageEdit.FormAddEmail.Form.SetHorizontal(false)
	groupsPage.PageEdit.FormAddEmail.Form.AddInputField("", "", 0, nil, nil)
	groupsPage.PageEdit.FormAddEmail.InputEmail = groupsPage.PageEdit.FormAddEmail.Form.GetFormItem(0).(*tview.InputField)

	//? Добавить e-mail
	groupsPage.PageEdit.ButtonAddEmail = tview.NewButton("Добавить e-mail")

	//? Удалить e-mail
	groupsPage.PageEdit.ButtonDeleteEmail = tview.NewButton("Удалить e-mail")

	//? Назад
	groupsPage.PageEdit.ButtonEditExit = tview.NewButton("Назад").SetSelectedFunc(func() {
		groupsPage.PagesSelEd.SwitchToPage("select_page")
		app.SetInputCapture(groupsPage.PageSelect.InputCapture)
		app.SetFocus(groupsPage.ListGroups)
	})

	//? EditButtonsGrid
	groupsPage.PageEdit.Grid = tview.NewGrid().
		SetRows(3, 1, 1, 1).
		SetColumns(0).
		SetGap(1, 0).
		AddItem(groupsPage.PageEdit.FormAddEmail.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageEdit.ButtonAddEmail, 1, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageEdit.ButtonDeleteEmail, 2, 0, 1, 1, 0, 0, true).
		AddItem(groupsPage.PageEdit.ButtonEditExit, 3, 0, 1, 1, 0, 0, true)

	//? emails list
	groupsPage.ListEmails = tview.NewList()
	groupsPage.ListEmails.SetBorder(true)
	groupsPage.ListEmails.SetHighlightFullLine(true)
	groupsPage.ListEmails.SetTitle(" Список допущенных e-mail ")

	//? InputCapture edit page
	groupsPage.PageEdit.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case groupsPage.ListEmails:
				app.SetFocus(groupsPage.PageEdit.FormAddEmail.InputEmail)
			case groupsPage.PageEdit.FormAddEmail.InputEmail:
				app.SetFocus(groupsPage.PageEdit.ButtonAddEmail)
			case groupsPage.PageEdit.ButtonAddEmail:
				app.SetFocus(groupsPage.PageEdit.ButtonDeleteEmail)
			case groupsPage.PageEdit.ButtonDeleteEmail:
				app.SetFocus(groupsPage.PageEdit.ButtonEditExit)
			case groupsPage.PageEdit.ButtonEditExit:
				app.SetFocus(groupsPage.ListEmails)
			}
			return nil
		}
		return event
	}

	//? Main grid
	groupsPage.GridMain = tview.NewGrid().
		SetRows(0).
		SetColumns(0, 30, 20, 30, 0).
		SetGap(1, 1).
		AddItem(groupsPage.ListGroups, 0, 1, 1, 1, 0, 0, true).
		AddItem(groupsPage.PagesSelEd, 0, 2, 1, 1, 0, 0, true).
		AddItem(groupsPage.ListEmails, 0, 3, 1, 1, 0, 0, true)

	// app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	// 	focus := app.GetFocus()
	// 	if focus != groupsPage.ButtonDeleteEmail {
	// 		return event, action

	// 	}
	// 	return nil, 0
	// })

	//? Adding pages
	groupsPage.PagesSelEd.AddPage("select_page", groupsPage.PageSelect.Grid, true, true)
	groupsPage.PagesSelEd.AddPage("edit_page", groupsPage.PageEdit.Grid, true, true)

	//! Adding Pages

	pages.AddPage("groups_page", groupsPage.GridMain, true, true)
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
