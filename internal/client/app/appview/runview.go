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

	//? pages
	loginPage := LoginPage{}
	registerPage := RegisterPage{}
	groupsPage := GroupsPage{}
	mainPage := PageMain{}

	// app
	mainPage.App = tview.NewApplication()

	// page container
	mainPage.Pages = tview.NewPages()

	//? left message box
	mainPage.MessageBoxL = tview.NewTextView()
	mainPage.MessageBoxL.SetDynamicColors(true)
	mainPage.MessageBoxL.SetTextAlign(tview.AlignLeft)
	mainPage.MessageBoxL.SetBorder(true).SetBorderColor(tcell.ColorRed)
	mainPage.MessageBoxL.SetTitle(" Сообщения ")

	//? right message box
	mainPage.MessageBoxR = tview.NewTextView()
	mainPage.MessageBoxR.SetDynamicColors(true)
	mainPage.MessageBoxR.SetTextAlign(tview.AlignLeft)
	mainPage.MessageBoxR.SetBorder(true).SetBorderColor(tcell.ColorRed)
	mainPage.MessageBoxR.SetTitle(" Дополнительная информация ")

	//? status bar
	mainPage.StatusBar.Table = tview.NewTable().SetBorders(true).SetBordersColor(tcell.ColorRed)
	mainPage.StatusBar.Table.SetCell(0, 0, tview.NewTableCell("Тип соединения").SetAlign(tview.AlignCenter).SetExpansion(1))
	mainPage.StatusBar.Table.SetCell(0, 1, tview.NewTableCell("Соединение с сервером").SetAlign(tview.AlignCenter).SetExpansion(1))
	mainPage.StatusBar.Table.SetCell(0, 2, tview.NewTableCell("Соединение с БД сервера").SetAlign(tview.AlignCenter).SetExpansion(1))
	mainPage.StatusBar.Table.SetCell(0, 3, tview.NewTableCell("Статус операции").SetAlign(tview.AlignCenter).SetExpansion(1))
	mainPage.StatusBar.Table.SetCell(0, 4, tview.NewTableCell("Активный аккаунт").SetAlign(tview.AlignCenter).SetExpansion(1))
	mainPage.StatusBar.CellTypeConnect = tview.NewTableCell("[green]REST API").SetAlign(tview.AlignCenter).SetExpansion(1)
	mainPage.StatusBar.Table.SetCell(1, 0, mainPage.StatusBar.CellTypeConnect)
	mainPage.StatusBar.CellServerConnect = tview.NewTableCell("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	mainPage.StatusBar.Table.SetCell(1, 1, mainPage.StatusBar.CellServerConnect)
	mainPage.StatusBar.CellServerDBConnect = tview.NewTableCell("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	mainPage.StatusBar.Table.SetCell(1, 2, mainPage.StatusBar.CellServerDBConnect)
	mainPage.StatusBar.CellResponseStatus = tview.NewTableCell("").SetAlign(tview.AlignCenter).SetExpansion(1)
	mainPage.StatusBar.Table.SetCell(1, 3, mainPage.StatusBar.CellResponseStatus)
	mainPage.StatusBar.CellActiveAccount = tview.NewTableCell("").SetAlign(tview.AlignCenter).SetExpansion(1)
	mainPage.StatusBar.Table.SetCell(1, 4, mainPage.StatusBar.CellActiveAccount)

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
			mainPage.App.QueueUpdateDraw(func() {
				mainPage.StatusBar.CellServerConnect.SetText(info)
			})
		}
	}()
	go func() {
		for info := range updateCellServerDBConnectChan {
			mainPage.App.QueueUpdateDraw(func() {
				mainPage.StatusBar.CellServerDBConnect.SetText(info)
			})
		}
	}()

	//? main grid
	mainGrid := tview.NewGrid()
	mainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	mainGrid.SetRows(0, 8, 5)
	mainGrid.AddItem(mainPage.Pages, 0, 0, 1, 2, 0, 0, true)
	mainGrid.AddItem(mainPage.MessageBoxL, 1, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(mainPage.MessageBoxR, 1, 1, 1, 1, 0, 0, true)
	mainGrid.AddItem(mainPage.StatusBar.Table, 2, 0, 1, 2, 0, 0, false)

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
		mainPage.Pages.SwitchToPage("groups_page")
		groupsPage.PagesSelEd.SwitchToPage("select_page")
		mainPage.App.SetInputCapture(groupsPage.PageSelect.InputCapture)
		mainPage.App.SetFocus(groupsPage.ListGroups)
	})

	//? Зарегистрироваться
	loginPage.ButtonRegister = tview.NewButton("Зарегистрироваться").SetSelectedFunc(func() {
		// switch
		mainPage.Pages.SwitchToPage("register_page")
		// focus
		mainPage.App.SetInputCapture(registerPage.InputCapture)
		mainPage.App.SetFocus(registerPage.Form.InputEmail)
		// clear
		registerPage.Form.InputEmail.SetText("")
		registerPage.Form.InputPassword.SetText("")
		registerPage.Form.InputRPassword.SetText("")
		// messageBox
		mainPage.MessageBoxL.SetText("Введите e-mail и пароль. Повторите ввод пароля. Нажмите кнопку [blue]\"Зарегистрироваться\".\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
	})

	//? Выход
	loginPage.ButtonExit = tview.NewButton("Выход").SetSelectedFunc(func() {
		mainPage.App.Stop()
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
			currentFocus := mainPage.App.GetFocus()
			switch currentFocus {
			case loginPage.Form.InputEmail:
				mainPage.App.SetFocus(loginPage.Form.InputPassword)
			case loginPage.Form.InputPassword:
				mainPage.App.SetFocus(loginPage.ButtonLogin)
			case loginPage.ButtonLogin:
				mainPage.App.SetFocus(loginPage.ButtonRegister)
			case loginPage.ButtonRegister:
				mainPage.App.SetFocus(loginPage.ButtonExit)
			case loginPage.ButtonExit:
				mainPage.App.SetFocus(loginPage.Form.InputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := mainPage.App.GetFocus()
		// 	switch currentFocus {
		// 	case formLoginPage.GetFormItem(1):
		// 		mainPage.App.SetFocus(buttonLoginLoginPage)
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
					mainPage.MessageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
					mainPage.MessageBoxL.SetText("[red]Возникла критическая ошибка.")
				} else {
					if registerAccountResp.HTTPStatusCode != 200 {
						mainPage.MessageBoxR.SetText(fmt.Sprintf("[red]%s", registerAccountResp.ServerError))
						switch {
						// invalid e-mail
						case strings.Contains(registerAccountResp.ServerError, customerrors.ErrAccountValidateEmail422.Error()):
							mainPage.MessageBoxL.SetText("[red]Неверно введён e-mail!")
							mainPage.App.SetFocus(registerPage.Form.InputEmail)
						// invalid password
						case strings.Contains(registerAccountResp.ServerError, customerrors.ErrAccountValidatePassword422.Error()):
							mainPage.MessageBoxL.SetText("[red]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов!")
							registerPage.Form.InputPassword.SetText("")
							registerPage.Form.InputRPassword.SetText("")
							mainPage.App.SetFocus(registerPage.Form.InputPassword)
						// e-mail is already busy
						case strings.Contains(registerAccountResp.ServerError, customerrors.ErrDBBusyEmail409.Error()):
							mainPage.MessageBoxL.SetText(fmt.Sprintf("[red]Пользователь с e-mail %s уже зарегестрирован!", email))
							mainPage.App.SetFocus(registerPage.Form.InputEmail)
						// other errors
						default:
							mainPage.MessageBoxL.SetText("[red]Возникла ошибка.")
						}
					} else {
						mainPage.MessageBoxR.Clear()
						mainPage.MessageBoxL.SetText(fmt.Sprintf("[green]Вы успешно зарегистрировались. Ваш ID: %s\n[white]Войдите в систему, используя свой e-mail и пароль.", registerAccountResp.AccountID))
					}
				}
			} else {
				mainPage.MessageBoxL.Clear()
				mainPage.MessageBoxL.SetText("[red]Пароли не совпадают! Повторите ввод.\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
				mainPage.MessageBoxR.Clear()
				registerPage.Form.InputPassword.SetText("")
				registerPage.Form.InputRPassword.SetText("")
				mainPage.App.SetFocus(registerPage.Form.InputPassword)
			}
		})

	//? Назад
	registerPage.ButtonExit = tview.NewButton("Назад").SetSelectedFunc(func() {
		// switch
		mainPage.Pages.SwitchToPage("login_page")
		// focus
		mainPage.App.SetInputCapture(loginPage.InputCapture)
		mainPage.App.SetFocus(loginPage.Form.InputEmail)
		// messageBox
		mainPage.MessageBoxL.Clear()
		mainPage.MessageBoxR.Clear()
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
			currentFocus := mainPage.App.GetFocus()
			switch currentFocus {
			case registerPage.Form.InputEmail:
				mainPage.App.SetFocus(registerPage.Form.InputPassword)
			case registerPage.Form.InputPassword:
				mainPage.App.SetFocus(registerPage.Form.InputRPassword)
			case registerPage.Form.InputRPassword:
				mainPage.App.SetFocus(registerPage.ButtonRegister)
			case registerPage.ButtonRegister:
				mainPage.App.SetFocus(registerPage.ButtonExit)
			case registerPage.ButtonExit:
				mainPage.App.SetFocus(registerPage.Form.InputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := mainPage.App.GetFocus()
		// 	switch currentFocus {
		// 	case registerPage.Form.Form.GetFormItem(2):
		// 		mainPage.App.SetFocus(registerPage.ButtonRegister)
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
		mainPage.MessageBoxL.SetText(mainText + secondaryText + string(shortcut))
	})

	//? "select_page"

	//? Выбрать группу
	groupsPage.PageSelect.ButtonSelect = tview.NewButton("Выбрать группу")

	//? Создать группу
	groupsPage.PageSelect.ButtonNew = tview.NewButton("Создать группу")

	//? Настроить группу
	groupsPage.PageSelect.ButtonSettings = tview.NewButton("Настроить группу").SetSelectedFunc(func() {
		groupsPage.PagesSelEd.SwitchToPage("edit_page")
		mainPage.App.SetInputCapture(groupsPage.PageEdit.InputCapture)
		mainPage.App.SetFocus(groupsPage.ListEmails)
		groupsPage.PageEdit.FormAddEmail.InputEmail.SetText("")
	})

	//? Удалить группу
	groupsPage.PageSelect.ButtonDelete = tview.NewButton("Удалить группу")

	//? Выйти из аккаунта
	groupsPage.PageSelect.ButtonLogout = tview.NewButton("Выйти из аккаунта").SetSelectedFunc(func() {
		// switch
		mainPage.Pages.SwitchToPage("login_page")
		// focus
		mainPage.App.SetInputCapture(loginPage.InputCapture)
		mainPage.App.SetFocus(loginPage.Form.InputEmail)
		// messageBox
		mainPage.MessageBoxL.Clear()
		mainPage.MessageBoxR.Clear()
	})

	//? Выход
	groupsPage.PageSelect.ButtonExit = tview.NewButton("Выход").SetSelectedFunc(func() {
		mainPage.App.Stop()
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
			currentFocus := mainPage.App.GetFocus()
			switch currentFocus {
			case groupsPage.ListGroups:
				mainPage.App.SetFocus(groupsPage.PageSelect.ButtonSelect)
			case groupsPage.PageSelect.ButtonSelect:
				mainPage.App.SetFocus(groupsPage.PageSelect.ButtonNew)
			case groupsPage.PageSelect.ButtonNew:
				mainPage.App.SetFocus(groupsPage.PageSelect.ButtonSettings)
			case groupsPage.PageSelect.ButtonSettings:
				mainPage.App.SetFocus(groupsPage.PageSelect.ButtonDelete)
			case groupsPage.PageSelect.ButtonDelete:
				mainPage.App.SetFocus(groupsPage.PageSelect.ButtonLogout)
			case groupsPage.PageSelect.ButtonLogout:
				mainPage.App.SetFocus(groupsPage.PageSelect.ButtonExit)
			case groupsPage.PageSelect.ButtonExit:
				mainPage.App.SetFocus(groupsPage.ListGroups)
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
		mainPage.App.SetInputCapture(groupsPage.PageSelect.InputCapture)
		mainPage.App.SetFocus(groupsPage.ListGroups)
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
			currentFocus := mainPage.App.GetFocus()
			switch currentFocus {
			case groupsPage.ListEmails:
				mainPage.App.SetFocus(groupsPage.PageEdit.FormAddEmail.InputEmail)
			case groupsPage.PageEdit.FormAddEmail.InputEmail:
				mainPage.App.SetFocus(groupsPage.PageEdit.ButtonAddEmail)
			case groupsPage.PageEdit.ButtonAddEmail:
				mainPage.App.SetFocus(groupsPage.PageEdit.ButtonDeleteEmail)
			case groupsPage.PageEdit.ButtonDeleteEmail:
				mainPage.App.SetFocus(groupsPage.PageEdit.ButtonEditExit)
			case groupsPage.PageEdit.ButtonEditExit:
				mainPage.App.SetFocus(groupsPage.ListEmails)
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

	// mainPage.App.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	// 	focus := mainPage.App.GetFocus()
	// 	if focus != groupsPage.ButtonDeleteEmail {
	// 		return event, action

	// 	}
	// 	return nil, 0
	// })

	//? Adding pages
	groupsPage.PagesSelEd.AddPage("select_page", groupsPage.PageSelect.Grid, true, true)
	groupsPage.PagesSelEd.AddPage("edit_page", groupsPage.PageEdit.Grid, true, true)

	//! Adding Pages

	mainPage.Pages.AddPage("groups_page", groupsPage.GridMain, true, true)
	mainPage.Pages.AddPage("register_page", registerPage.MainGrid, true, true)
	mainPage.Pages.AddAndSwitchToPage("login_page", loginPage.MainGrid, true)
	mainPage.App.SetInputCapture(loginPage.InputCapture)
	mainPage.App.SetFocus(loginPage.Form.InputEmail)

	//! Launching the app
	mainPage.App.SetRoot(mainGrid, true)

	if err := mainPage.App.Run(); err != nil {
		return fmt.Errorf("%s: %s: %w: %w", customerrors.ClientAppViewErr, action, customerrors.ErrRunAppView, err)
	}

	return nil
}
