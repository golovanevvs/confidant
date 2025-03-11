package appview

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

type FormPageRegister struct {
	Form           *tview.Form
	InputEmail     *tview.InputField
	InputPassword  *tview.InputField
	InputRPassword *tview.InputField
}

type PageRegister struct {
	InputCapture   func(event *tcell.EventKey) *tcell.EventKey
	Form           FormPageRegister
	ButtonRegister *tview.Button
	ButtonExit     *tview.Button
	Grid           *tview.Grid
	MainGrid       *tview.Grid
}

func (av *AppView) VRegister() {
	//? Form
	av.v.pageRegister.Form.Form.SetHorizontal(false)
	av.v.pageRegister.Form.Form.AddInputField("E-mail:", "", 0, nil, nil)
	av.v.pageRegister.Form.InputEmail = av.v.pageRegister.Form.Form.GetFormItem(0).(*tview.InputField)
	av.v.pageRegister.Form.Form.AddPasswordField("Пароль:", "", 0, '*', nil)
	av.v.pageRegister.Form.InputPassword = av.v.pageRegister.Form.Form.GetFormItem(1).(*tview.InputField)
	av.v.pageRegister.Form.Form.AddPasswordField("Повторите:", "", 0, '*', nil)
	av.v.pageRegister.Form.InputRPassword = av.v.pageRegister.Form.Form.GetFormItem(2).(*tview.InputField)

	//? Зарегистрироваться
	av.v.pageRegister.ButtonRegister.
		SetSelectedFunc(func() {
			pass1 := av.v.pageRegister.Form.InputPassword.GetText()
			pass2 := av.v.pageRegister.Form.InputRPassword.GetText()
			if pass1 == pass2 {
				email := av.v.pageRegister.Form.InputEmail.GetText()
				password := av.v.pageRegister.Form.InputPassword.GetText()
				registerAccountResp, err := av.sv.RegisterAccount(email, password)
				if err != nil {
					av.v.pageMain.MessageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))
					av.v.pageMain.MessageBoxL.SetText("[red]Возникла критическая ошибка.")
				} else {
					av.v.pageMain.MessageBoxR.SetText(fmt.Sprintf("[red]%s", registerAccountResp.Error))
					// setting status
					if registerAccountResp.HTTPStatusCode == 200 {
						av.v.pageMain.StatusBar.CellResponseStatus.SetText(fmt.Sprintf("[green]%s", registerAccountResp.HTTPStatus))
					} else {
						av.v.pageMain.StatusBar.CellResponseStatus.SetText(fmt.Sprintf("[yellow]%s", registerAccountResp.HTTPStatus))
					}
					switch {
					case registerAccountResp.HTTPStatusCode != 200:
						switch {
						// invalid e-mail
						case strings.Contains(registerAccountResp.Error, customerrors.ErrAccountValidateEmail422.Error()):
							av.v.pageMain.MessageBoxL.SetText("[red]Неверно введён e-mail!")
							av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputEmail)
						// invalid password
						case strings.Contains(registerAccountResp.Error, customerrors.ErrAccountValidatePassword422.Error()):
							av.v.pageMain.MessageBoxL.SetText("[red]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов!")
							av.v.pageRegister.Form.InputPassword.SetText("")
							av.v.pageRegister.Form.InputRPassword.SetText("")
							av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputPassword)
						// e-mail is already busy
						case strings.Contains(registerAccountResp.Error, customerrors.ErrDBBusyEmail409.Error()):
							av.v.pageMain.MessageBoxL.SetText(fmt.Sprintf("[red]Пользователь с e-mail %s уже зарегистрирован!", email))
							av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputEmail)
						// other errors
						default:
							av.v.pageMain.MessageBoxL.SetText("[red]Возникла ошибка.")
						}

					// OK
					case registerAccountResp.Error == "":
						av.v.pageMain.StatusBar.CellActiveAccount.SetText(fmt.Sprintf("[green]%s", email))
						av.v.pageMain.MessageBoxL.Clear()
						av.v.pageMain.MessageBoxR.Clear()
						av.v.pageMain.Pages.SwitchToPage("groups_page")
						av.v.pageGroups.PagesSelEd.SwitchToPage("select_page")
						av.v.pageMain.App.SetInputCapture(av.v.pageGroups.PageSelect.InputCapture)
						av.v.pageMain.App.SetFocus(av.v.pageGroups.ListGroups)
					// the response does not contain the "Authorization" header
					case strings.Contains(registerAccountResp.Error, customerrors.ErrAuthHeader.Error()):
						av.v.pageMain.MessageBoxL.SetText("[red]Ответ не содержит заголовок \"Authorization\"!")
					// the "Authorization" header does not contain "Bearer"
					case strings.Contains(registerAccountResp.Error, customerrors.ErrBearer.Error()):
						av.v.pageMain.MessageBoxL.SetText("[red]Заголовок \"Authorization\" не содержит \"Bearer\"!")
					// the "Authorization" header does not contain a access token
					case strings.Contains(registerAccountResp.Error, customerrors.ErrAccessToken.Error()):
						av.v.pageMain.MessageBoxL.SetText("[red]Заголовок \"Authorization\" не содержит access токен!")
					// the response does not contain the "Refresh-Token" header
					case strings.Contains(registerAccountResp.Error, customerrors.ErrRefreshToken.Error()):
						av.v.pageMain.MessageBoxL.SetText("[red]Заголовок \"Authorization\" не содержит refresh токен!")
					// unknown error
					default:
						av.v.pageMain.MessageBoxL.SetText("Неизвестная ошибка!")
					}
				}
			} else {
				av.v.pageMain.MessageBoxL.Clear()
				av.v.pageMain.MessageBoxL.SetText("[red]Пароли не совпадают! Повторите ввод.\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
				av.v.pageMain.MessageBoxR.Clear()
				av.v.pageRegister.Form.InputPassword.SetText("")
				av.v.pageRegister.Form.InputRPassword.SetText("")
				av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputPassword)
			}
		})

	//? Назад
	av.v.pageRegister.ButtonExit.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.Pages.SwitchToPage("login_page")
		// focus
		av.v.pageMain.App.SetInputCapture(av.v.pageLogin.InputCapture)
		av.v.pageMain.App.SetFocus(av.v.pageLogin.Form.InputEmail)
		// messageBox
		av.v.pageMain.MessageBoxL.Clear()
		av.v.pageMain.MessageBoxR.Clear()
		av.v.pageMain.StatusBar.CellResponseStatus.SetText("")
	})

	//? form grid
	av.v.pageRegister.Grid.
		SetRows(8, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(av.v.pageRegister.Form.Form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageRegister.ButtonRegister, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageRegister.ButtonExit, 2, 0, 1, 1, 0, 0, true)

	//? main grid
	av.v.pageRegister.MainGrid.
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(av.v.pageRegister.Grid, 1, 1, 1, 1, 0, 0, true)

	//? InputCapture
	av.v.pageRegister.InputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageMain.App.GetFocus()
			switch currentFocus {
			case av.v.pageRegister.Form.InputEmail:
				av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputPassword)
			case av.v.pageRegister.Form.InputPassword:
				av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputRPassword)
			case av.v.pageRegister.Form.InputRPassword:
				av.v.pageMain.App.SetFocus(av.v.pageRegister.ButtonRegister)
			case av.v.pageRegister.ButtonRegister:
				av.v.pageMain.App.SetFocus(av.v.pageRegister.ButtonExit)
			case av.v.pageRegister.ButtonExit:
				av.v.pageMain.App.SetFocus(av.v.pageRegister.Form.InputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := av.v.pageMain.App.GetFocus()
		// 	switch currentFocus {
		// 	case av.v.pageRegister.Form.Form.GetFormItem(2):
		// 		av.v.pageMain.App.SetFocus(av.v.pageRegister.ButtonRegister)
		// 	}
		// }

		return event
	}
}
