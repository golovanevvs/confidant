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

	//! LOGIN PAGE
	formLoginPage = tview.NewForm()
	formLoginPage.SetHorizontal(false)
	formLoginPage.AddInputField("E-mail:", "", 0, nil, nil)
	formLoginPage.AddPasswordField("Пароль:", "", 0, '*', nil)
	//! Войти
	buttonLoginLoginPage := tview.NewButton("Войти")
	//! Зарегистрироваться
	buttonRegisterLoginPage := tview.NewButton("Зарегистрироваться").SetSelectedFunc(func() {
		pages.SwitchToPage("register_page")
		app.SetInputCapture(inputCaptureRegisterPage)
		app.SetFocus(formRegisterPage.GetFormItem(0))
		for i := 0; i < 3; i++ {
			if inputField, ok := formRegisterPage.GetFormItem(i).(*tview.InputField); ok {
				inputField.SetText("")
			}
		}
		messageBoxL.SetText("Введите e-mail и пароль. Повторите ввод пароля. Нажмите кнопку [blue]\"Зарегистрироваться\".\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
	})
	//! Выход
	buttonExitLoginPage := tview.NewButton("Выход").SetSelectedFunc(func() {
		app.Stop()
	})

	gridLoginPage := tview.NewGrid().
		SetRows(5, 3, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(formLoginPage, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonLoginLoginPage, 1, 0, 1, 1, 1, 1, false).
		AddItem(buttonRegisterLoginPage, 2, 0, 1, 1, 0, 0, false).
		AddItem(buttonExitLoginPage, 3, 0, 1, 1, 0, 0, false)

	mainGridLoginPage := tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(gridLoginPage, 1, 1, 1, 1, 0, 0, true)

	inputCaptureLoginPage = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case formLoginPage.GetFormItem(0):
				app.SetFocus(formLoginPage.GetFormItem(1))
			case formLoginPage.GetFormItem(1):
				app.SetFocus(buttonLoginLoginPage)
			case buttonLoginLoginPage:
				app.SetFocus(buttonRegisterLoginPage)
			case buttonRegisterLoginPage:
				app.SetFocus(buttonExitLoginPage)
			case buttonExitLoginPage:
				app.SetFocus(formLoginPage.GetFormItem(0))
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
	formRegisterPage.AddPasswordField("Пароль:", "", 0, '*', nil)
	formRegisterPage.AddPasswordField("Повторите:", "", 0, '*', nil)

	buttonRegisterRegisterPage := tview.NewButton("Зарегистрироваться").
		SetSelectedFunc(func() {
			//! Зарегистрироваться
			pass1 := formRegisterPage.GetFormItem(1).(*tview.InputField).GetText()
			pass2 := formRegisterPage.GetFormItem(2).(*tview.InputField).GetText()
			if pass1 != pass2 {
				messageBoxL.SetText("[red]Пароли не совпадают! Повторите ввод.\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
				messageBoxR.Clear()
				if inputField, ok := formRegisterPage.GetFormItem(1).(*tview.InputField); ok {
					inputField.SetText("")
					app.SetFocus(inputField)
				}
				if inputField, ok := formRegisterPage.GetFormItem(2).(*tview.InputField); ok {
					inputField.SetText("")
				}
			} else {
				email := formRegisterPage.GetFormItem(0).(*tview.InputField).GetText()
				password := formRegisterPage.GetFormItem(1).(*tview.InputField).GetText()
				accountID, err := av.sv.RegisterAccount(email, password)
				if err != nil {
					messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
					switch {
					case strings.Contains(err.Error(), customerrors.ErrAccountValidateEmail422.Error()):
						messageBoxL.SetText("[red]Неверно введён e-mail!")
						if inputField, ok := formRegisterPage.GetFormItem(0).(*tview.InputField); ok {
							app.SetFocus(inputField)
						}
					case strings.Contains(err.Error(), customerrors.ErrAccountValidatePassword422.Error()):
						messageBoxL.SetText("[red]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов!")
						if inputField, ok := formRegisterPage.GetFormItem(1).(*tview.InputField); ok {
							inputField.SetText("")
							app.SetFocus(inputField)
						}
						if inputField, ok := formRegisterPage.GetFormItem(2).(*tview.InputField); ok {
							inputField.SetText("")
						}
					case strings.Contains(err.Error(), customerrors.ErrDBBusyEmail409.Error()):
						messageBoxL.SetText(fmt.Sprintf("[red]Пользователь с e-mail %s уже зарегестрирован!", email))
						if inputField, ok := formRegisterPage.GetFormItem(0).(*tview.InputField); ok {
							app.SetFocus(inputField)
						}
					default:
						messageBoxL.SetText("[red]Возникла ошибка.")
					}
				} else {
					messageBoxL.Clear()
					messageBoxR.Clear()
					messageBoxL.SetText(fmt.Sprintf("[green]Вы успешно зарегистрировались. Ваш ID: %d\n[white]Войдите в систему, используя свой e-mail и пароль.", accountID))
				}
			}
		})
		//! Назад
	buttonExitRegisterPage := tview.NewButton("Назад").SetSelectedFunc(func() {
		pages.AddAndSwitchToPage("login_page", mainGridLoginPage, true)
		app.SetInputCapture(inputCaptureLoginPage)
		app.SetFocus(formLoginPage.GetFormItem(0))
		messageBoxL.SetText("[green]Добро пожаловать в систему хранения конфиденциальной информации [white]CON[blue]FID[red]ANT")
		messageBoxR.Clear()
	})

	formGrid := tview.NewGrid().
		SetRows(8, 3, 3).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(formRegisterPage, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonRegisterRegisterPage, 1, 0, 1, 1, 0, 0, true).
		AddItem(buttonExitRegisterPage, 2, 0, 1, 1, 0, 0, true)

	mainGridRegisterPage := tview.NewGrid().
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(formGrid, 1, 1, 1, 1, 0, 0, true)

	inputCaptureRegisterPage = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case formRegisterPage.GetFormItem(0):
				app.SetFocus(formRegisterPage.GetFormItem(1))
			case formRegisterPage.GetFormItem(1):
				app.SetFocus(formRegisterPage.GetFormItem(2))
			case formRegisterPage.GetFormItem(2):
				app.SetFocus(buttonRegisterRegisterPage)
			case buttonRegisterRegisterPage:
				app.SetFocus(buttonExitRegisterPage)
			case buttonExitRegisterPage:
				app.SetFocus(formRegisterPage.GetFormItem(0))
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
	app.SetFocus(formLoginPage.GetFormItem(0))

	//! Launching the app
	app.SetRoot(mainGrid, true)

	if err := app.Run(); err != nil {
		return fmt.Errorf("%s: %s: %w: %w", customerrors.ClientAppViewErr, action, customerrors.ErrRunAppView, err)
	}

	return nil
}
