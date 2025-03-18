package appview

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

type formPageLogin struct {
	form          *tview.Form
	inputEmail    *tview.InputField
	inputPassword *tview.InputField
}

type pageLogin struct {
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	form           formPageLogin
	buttonLogin    *tview.Button
	buttonRegister *tview.Button
	buttonExit     *tview.Button
	grid           *tview.Grid
	mainGrid       *tview.Grid
}

func newPageLogin() *pageLogin {
	return &pageLogin{
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		form: formPageLogin{
			form:          tview.NewForm(),
			inputEmail:    tview.NewInputField(),
			inputPassword: tview.NewInputField(),
		},
		buttonLogin:    tview.NewButton("Войти"),
		buttonRegister: tview.NewButton("Регистрация"),
		buttonExit:     tview.NewButton("Выход"),
		grid:           tview.NewGrid(),
		mainGrid:       tview.NewGrid(),
	}
}

func (av *appView) vLogin() {
	//! form
	av.v.pageLogin.form.form.SetHorizontal(false)
	av.v.pageLogin.form.form.AddInputField("E-mail:", "", 0, nil, nil)
	av.v.pageLogin.form.inputEmail = av.v.pageLogin.form.form.GetFormItem(0).(*tview.InputField)
	av.v.pageLogin.form.form.AddPasswordField("Пароль:", "", 0, '*', nil)
	av.v.pageLogin.form.inputPassword = av.v.pageLogin.form.form.GetFormItem(1).(*tview.InputField)

	//! "Войти"
	av.v.pageLogin.buttonLogin.SetSelectedFunc(func() {

		email := av.v.pageLogin.form.inputEmail.GetText()
		password := av.v.pageLogin.form.inputPassword.GetText()

		accountResp, err := av.sv.Login(context.Background(), email, password)

		// error
		if err != nil {
			av.v.pageMain.statusBar.cellResponseStatus.SetText("")
			if errors.Is(err, customerrors.ErrDBWrongPassword401) {
				av.v.pageMain.statusBar.cellResponseStatus.SetText("")
				av.v.pageMain.messageBoxL.SetText("[red]Введён неверный пароль!")
				av.v.pageMain.messageBoxR.SetText("")
				av.v.pageLogin.form.inputPassword.SetText("")
				av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputPassword)
			} else {
				av.v.pageMain.messageBoxL.SetText("[red]Возникла критическая ошибка.")
				av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
			}

			// no error
		} else {
			av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", accountResp.Error))

			// setting status
			if accountResp.HTTPStatusCode == 200 {
				av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[green]%s", accountResp.HTTPStatus))
			} else {
				av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[yellow]%s", accountResp.HTTPStatus))
				switch {
				case strings.Contains(accountResp.Error, customerrors.ErrDBEmailNotFound401.Error()):
					av.v.pageMain.messageBoxL.SetText(fmt.Sprintf("[red] e-mail %s не зарегистрирован!", email))
					av.v.pageLogin.form.inputPassword.SetText("")
					av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputEmail)
				case strings.Contains(accountResp.Error, customerrors.ErrDBWrongPassword401.Error()):
					av.v.pageMain.messageBoxL.SetText("[red]Введён неверный пароль!")
					av.v.pageLogin.form.inputPassword.SetText("")
					av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputPassword)
				default:
					av.v.pageMain.messageBoxL.SetText("[red]Возникла ошибка.")
				}
			}

			if accountResp.Error == "" {
				av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[green]%s", accountResp.HTTPStatus))
				av.v.pageMain.statusBar.cellActiveAccount.SetText(fmt.Sprintf("[green]%s", email))
				av.v.pageMain.messageBoxL.Clear()
				av.v.pageMain.messageBoxR.Clear()

				// switch
				av.v.pageMain.pages.SwitchToPage("groups_page")
				av.v.pageGroups.pagesSelEd.SwitchToPage("select_page")
				av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageSelect.inputCapture)
				av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)
			}

		}

	})

	//! "Регистрация"
	av.v.pageLogin.buttonRegister.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.pages.SwitchToPage("register_page")
		// focus
		av.v.pageApp.app.SetInputCapture(av.v.pageRegister.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputEmail)
		// clear
		av.v.pageRegister.form.inputEmail.SetText("")
		av.v.pageRegister.form.inputPassword.SetText("")
		av.v.pageRegister.form.inputRPassword.SetText("")
		// messageBox
		av.v.pageMain.messageBoxL.SetText("Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
	})

	//! "Выход"
	av.v.pageLogin.buttonExit.SetSelectedFunc(func() {
		av.v.pageApp.app.Stop()
	})

	//! form grid
	av.v.pageLogin.grid.
		SetRows(5, 1, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(av.v.pageLogin.form.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageLogin.buttonLogin, 1, 0, 1, 1, 1, 1, false).
		AddItem(av.v.pageLogin.buttonRegister, 2, 0, 1, 1, 0, 0, false).
		AddItem(av.v.pageLogin.buttonExit, 3, 0, 1, 1, 0, 0, false)

	//! main grid
	av.v.pageLogin.mainGrid.
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(av.v.pageLogin.grid, 1, 1, 1, 1, 0, 0, true)

	//! InputCapture
	av.v.pageLogin.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageLogin.form.inputEmail:
				av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputPassword)
			case av.v.pageLogin.form.inputPassword:
				av.v.pageApp.app.SetFocus(av.v.pageLogin.buttonLogin)
			case av.v.pageLogin.buttonLogin:
				av.v.pageApp.app.SetFocus(av.v.pageLogin.buttonRegister)
			case av.v.pageLogin.buttonRegister:
				av.v.pageApp.app.SetFocus(av.v.pageLogin.buttonExit)
			case av.v.pageLogin.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := av.v.pageApp.App.GetFocus()
		// 	switch currentFocus {
		// 	case formav.v.pageLogin.GetFormItem(1):
		// 		av.v.pageApp.App.SetFocus(buttonLoginav.v.pageLogin)
		// 	}
		// }

		return event
	}
}
