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

	var formRegisterPage, formLoginPage *tview.Form
	var inputCaptureLoginPage, inputCaptureRegisterPage func(event *tcell.EventKey) *tcell.EventKey

	// login page
	var inputEmailLoginPage, inputPasswordLoginPage *tview.InputField
	var buttonLoginLoginPage, buttonRegisterLoginPage, buttonExitLoginPage *tview.Button
	var gridLoginPage, mainGridLoginPage *tview.Grid

	// register page
	var inputEmailRegisterPage, inputPasswordRegisterPage, inputRPasswordRegisterPage *tview.InputField
	var buttonRegisterRegisterPage, buttonExitRegisterPage *tview.Button
	var gridRegisterPage, mainGridRegisterPage *tview.Grid

	//! LOGIN PAGE
	formLoginPage = tview.NewForm()
	formLoginPage.SetHorizontal(false)
	formLoginPage.AddInputField("E-mail:", "", 0, nil, nil)
	inputEmailLoginPage = formLoginPage.GetFormItem(0).(*tview.InputField)
	formLoginPage.AddPasswordField("Пароль:", "", 0, '*', nil)
	inputPasswordLoginPage = formLoginPage.GetFormItem(1).(*tview.InputField)
	//! Войти
	buttonLoginLoginPage = tview.NewButton("Войти")
	//! Зарегистрироваться
	buttonRegisterLoginPage = tview.NewButton("Зарегистрироваться").SetSelectedFunc(func() {
		// switch
		pages.SwitchToPage("register_page")
		// focus
		app.SetInputCapture(inputCaptureRegisterPage)
		app.SetFocus(inputEmailRegisterPage)
		// clear
		inputEmailRegisterPage.SetText("")
		inputPasswordRegisterPage.SetText("")
		inputRPasswordRegisterPage.SetText("")
		// messageBox
		messageBoxL.SetText("Введите e-mail и пароль. Повторите ввод пароля. Нажмите кнопку [blue]\"Зарегистрироваться\".\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
	})
	//! Выход
	buttonExitLoginPage = tview.NewButton("Выход").SetSelectedFunc(func() {
		app.Stop()
	})

	gridLoginPage = tview.NewGrid().
		SetRows(5, 3, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(formLoginPage, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonLoginLoginPage, 1, 0, 1, 1, 1, 1, false).
		AddItem(buttonRegisterLoginPage, 2, 0, 1, 1, 0, 0, false).
		AddItem(buttonExitLoginPage, 3, 0, 1, 1, 0, 0, false)

	mainGridLoginPage = tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(gridLoginPage, 1, 1, 1, 1, 0, 0, true)

	inputCaptureLoginPage = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case inputEmailLoginPage:
				app.SetFocus(inputPasswordLoginPage)
			case inputPasswordLoginPage:
				app.SetFocus(buttonLoginLoginPage)
			case buttonLoginLoginPage:
				app.SetFocus(buttonRegisterLoginPage)
			case buttonRegisterLoginPage:
				app.SetFocus(buttonExitLoginPage)
			case buttonExitLoginPage:
				app.SetFocus(inputEmailLoginPage)
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
	formRegisterPage = tview.NewForm()
	formRegisterPage.SetHorizontal(false)
	formRegisterPage.AddInputField("E-mail:", "", 0, nil, nil)
	inputEmailRegisterPage = formRegisterPage.GetFormItem(0).(*tview.InputField)
	formRegisterPage.AddPasswordField("Пароль:", "", 0, '*', nil)
	inputPasswordRegisterPage = formRegisterPage.GetFormItem(1).(*tview.InputField)
	formRegisterPage.AddPasswordField("Повторите:", "", 0, '*', nil)
	inputRPasswordRegisterPage = formRegisterPage.GetFormItem(2).(*tview.InputField)

	buttonRegisterRegisterPage = tview.NewButton("Зарегистрироваться").
		SetSelectedFunc(func() {
			//! Зарегистрироваться
			pass1 := inputPasswordRegisterPage.GetText()
			pass2 := inputRPasswordRegisterPage.GetText()
			if pass1 == pass2 {
				email := inputEmailRegisterPage.GetText()
				password := inputPasswordRegisterPage.GetText()
				registerAccountResp, err := av.sv.RegisterAccount(email, password)
				if err != nil {
					messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
					switch {
					// invalid e-mail
					case strings.Contains(err.Error(), customerrors.ErrAccountValidateEmail422.Error()):
						messageBoxL.SetText("[red]Неверно введён e-mail!")
						app.SetFocus(inputEmailRegisterPage)
					// invalid password
					case strings.Contains(err.Error(), customerrors.ErrAccountValidatePassword422.Error()):
						messageBoxL.SetText("[red]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов!")
						inputPasswordRegisterPage.SetText("")
						inputRPasswordRegisterPage.SetText("")
						app.SetFocus(inputPasswordRegisterPage)
					// e-mail is already busy
					case strings.Contains(err.Error(), customerrors.ErrDBBusyEmail409.Error()):
						messageBoxL.SetText(fmt.Sprintf("[red]Пользователь с e-mail %s уже зарегестрирован!", email))
						app.SetFocus(inputEmailRegisterPage)
					// other errors
					default:
						messageBoxL.SetText("[red]Возникла ошибка.")
					}
				} else {
					messageBoxL.Clear()
					messageBoxR.Clear()
					messageBoxL.SetText(fmt.Sprintf("[green]Вы успешно зарегистрировались. Ваш ID: %s\n[white]Войдите в систему, используя свой e-mail и пароль.", registerAccountResp.AccountID))
				}
			} else {
				messageBoxL.Clear()
				messageBoxL.SetText("[red]Пароли не совпадают! Повторите ввод.\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
				messageBoxR.Clear()
				inputPasswordRegisterPage.SetText("")
				inputRPasswordRegisterPage.SetText("")
				app.SetFocus(inputPasswordRegisterPage)
			}
		})
		//! Назад
	buttonExitRegisterPage = tview.NewButton("Назад").SetSelectedFunc(func() {
		// switch
		// pages.AddAndSwitchToPage("login_page", mainGridLoginPage, true)
		pages.SwitchToPage("login_page")
		// focus
		app.SetInputCapture(inputCaptureLoginPage)
		app.SetFocus(inputEmailLoginPage)
		// messageBox
		messageBoxL.SetText("[green]Добро пожаловать в систему хранения конфиденциальной информации [white]CON[blue]FID[red]ANT")
		messageBoxR.Clear()
	})

	gridRegisterPage = tview.NewGrid().
		SetRows(8, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(formRegisterPage, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonRegisterRegisterPage, 1, 0, 1, 1, 0, 0, true).
		AddItem(buttonExitRegisterPage, 2, 0, 1, 1, 0, 0, true)

	mainGridRegisterPage = tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(gridRegisterPage, 1, 1, 1, 1, 0, 0, true)

	inputCaptureRegisterPage = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case inputEmailRegisterPage:
				app.SetFocus(inputPasswordRegisterPage)
			case inputPasswordRegisterPage:
				app.SetFocus(inputRPasswordRegisterPage)
			case inputRPasswordRegisterPage:
				app.SetFocus(buttonRegisterRegisterPage)
			case buttonRegisterRegisterPage:
				app.SetFocus(buttonExitRegisterPage)
			case buttonExitRegisterPage:
				app.SetFocus(inputEmailRegisterPage)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := app.GetFocus()
		// 	switch currentFocus {
		// 	case formRegisterPage.GetFormItem(2):
		// 		app.SetFocus(buttonRegisterRegisterPage)
		// 	}
		// }

		return event
	}

	//! Adding Pages

	pages.AddPage("register_page", mainGridRegisterPage, true, true)
	pages.AddAndSwitchToPage("login_page", mainGridLoginPage, true)
	app.SetInputCapture(inputCaptureLoginPage)
	app.SetFocus(inputEmailLoginPage)

	//! Launching the app
	app.SetRoot(mainGrid, true)

	if err := app.Run(); err != nil {
		return fmt.Errorf("%s: %s: %w: %w", customerrors.ClientAppViewErr, action, customerrors.ErrRunAppView, err)
	}

	return nil
}
