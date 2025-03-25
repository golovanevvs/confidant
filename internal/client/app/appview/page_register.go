package appview

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

type formPageRegister struct {
	form           *tview.Form
	inputEmail     *tview.InputField
	inputPassword  *tview.InputField
	inputRPassword *tview.InputField
}

type pageRegister struct {
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	form           formPageRegister
	buttonRegister *tview.Button
	buttonExit     *tview.Button
	grid           *tview.Grid
	mainGrid       *tview.Grid
}

func newPageRegister() *pageRegister {
	return &pageRegister{
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		form: formPageRegister{
			form:           tview.NewForm(),
			inputEmail:     tview.NewInputField(),
			inputPassword:  tview.NewInputField(),
			inputRPassword: tview.NewInputField(),
		},
		buttonRegister: tview.NewButton("Зарегистрироваться"),
		buttonExit:     tview.NewButton("Назад"),
		grid:           tview.NewGrid(),
		mainGrid:       tview.NewGrid(),
	}
}

func (av *appView) vRegister() {
	//! form
	av.v.pageRegister.form.form.SetHorizontal(false)
	av.v.pageRegister.form.form.AddInputField("E-mail:", "", 0, nil, nil)
	av.v.pageRegister.form.inputEmail = av.v.pageRegister.form.form.GetFormItem(0).(*tview.InputField)
	av.v.pageRegister.form.form.AddPasswordField("Пароль:", "", 0, '*', nil)
	av.v.pageRegister.form.inputPassword = av.v.pageRegister.form.form.GetFormItem(1).(*tview.InputField)
	av.v.pageRegister.form.form.AddPasswordField("Повторите:", "", 0, '*', nil)
	av.v.pageRegister.form.inputRPassword = av.v.pageRegister.form.form.GetFormItem(2).(*tview.InputField)

	//! "Зарегистрироваться"
	av.v.pageRegister.buttonRegister.
		SetSelectedFunc(func() {
			pass1 := av.v.pageRegister.form.inputPassword.GetText()
			pass2 := av.v.pageRegister.form.inputRPassword.GetText()

			// pass1 == pass2
			if pass1 == pass2 {
				email := av.v.pageRegister.form.inputEmail.GetText()
				password := av.v.pageRegister.form.inputPassword.GetText()

				//? running service
				accountResp, err := av.sv.CreateAccount(context.Background(), email, password)

				// error
				if err != nil {
					av.v.pageMain.statusBar.cellResponseStatus.SetText("")
					av.v.pageMain.messageBoxL.SetText("[red]Ошибка.")
					av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", err.Error()))

					// no error
				} else {
					av.v.pageMain.messageBoxR.SetText(fmt.Sprintf("[red]%s", accountResp.Error))

					// setting status
					if accountResp.HTTPStatusCode == 200 {
						av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[green]%s", accountResp.HTTPStatus))
					} else {
						av.v.pageMain.statusBar.cellResponseStatus.SetText(fmt.Sprintf("[yellow]%s", accountResp.HTTPStatus))
					}

					switch {

					// status != 200

					case accountResp.HTTPStatusCode != 200:
						switch {
						// invalid e-mail
						case strings.Contains(accountResp.Error, customerrors.ErrAccountValidateEmail422.Error()):
							av.v.pageMain.messageBoxL.SetText("[red]Неверно введён e-mail!")
							av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputEmail)
						// invalid password
						case strings.Contains(accountResp.Error, customerrors.ErrAccountValidatePassword422.Error()):
							av.v.pageMain.messageBoxL.SetText("[red]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов!")
							av.v.pageRegister.form.inputPassword.SetText("")
							av.v.pageRegister.form.inputRPassword.SetText("")
							av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputPassword)
						// e-mail is already busy
						case strings.Contains(accountResp.Error, customerrors.ErrDBBusyEmail409.Error()):
							av.v.pageMain.messageBoxL.SetText(fmt.Sprintf("[red]Пользователь с e-mail %s уже зарегистрирован!", email))
							av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputEmail)
						// other errors
						default:
							av.v.pageMain.messageBoxL.SetText("[red]Возникла ошибка.")
						}

					// status == 200

					//? OK
					case accountResp.Error == "":
						av.account.ID = accountResp.AccountID
						av.account.Email = email
						av.v.pageMain.statusBar.cellActiveAccount.SetText(fmt.Sprintf("[green]%s", email))
						av.v.pageMain.messageBoxL.Clear()
						av.v.pageMain.messageBoxR.Clear()

						// updating groups list
						av.v.pageGroups.listGroups.Clear()
						if len(av.groups) > 0 {
							for _, group := range av.groups {
								av.v.pageGroups.listGroups.AddItem(group.Title, "", 0, nil)
							}

							// updating e-mails
							av.v.pageGroups.listEmails.Clear()
							for _, email := range av.groups[0].Emails {
								av.v.pageGroups.listEmails.AddItem(email, "", 0, nil)
							}
						}

						av.v.pageMain.pages.SwitchToPage("groups_page")
						av.v.pageGroups.pages.SwitchToPage("select_page")
						av.v.pageApp.app.SetInputCapture(av.v.pageGroups.pageGroupsSelect.inputCapture)
						av.v.pageApp.app.SetFocus(av.v.pageGroups.listGroups)

						av.accessToken = accountResp.AccessTokenString

					// the response does not contain the "Authorization" header
					case strings.Contains(accountResp.Error, customerrors.ErrAuthHeaderResp.Error()):
						av.v.pageMain.messageBoxL.SetText("[red]Ответ не содержит заголовок \"Authorization\"!")
					// the "Authorization" header does not contain "Bearer"
					case strings.Contains(accountResp.Error, customerrors.ErrBearer.Error()):
						av.v.pageMain.messageBoxL.SetText("[red]Заголовок \"Authorization\" не содержит \"Bearer\"!")
					// the "Authorization" header does not contain a access token
					case strings.Contains(accountResp.Error, customerrors.ErrAccessToken.Error()):
						av.v.pageMain.messageBoxL.SetText("[red]Заголовок \"Authorization\" не содержит access токен!")
					// the response does not contain the "Refresh-Token" header
					case strings.Contains(accountResp.Error, customerrors.ErrRefreshToken.Error()):
						av.v.pageMain.messageBoxL.SetText("[red]Заголовок \"Authorization\" не содержит refresh токен!")
					// unknown error
					default:
						av.v.pageMain.messageBoxL.SetText("Неизвестная ошибка!")
					}
				}

				// pass1 != pass2
			} else {
				av.v.pageMain.messageBoxL.Clear()
				av.v.pageMain.messageBoxL.SetText("[red]Пароли не совпадают! Повторите ввод.\n[white]Пароль должен содержать минимум 8 символов, состоять из заглавных и строчных букв латинского алфавита, цифр и символов.")
				av.v.pageMain.messageBoxR.Clear()
				av.v.pageRegister.form.inputPassword.SetText("")
				av.v.pageRegister.form.inputRPassword.SetText("")
				av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputPassword)
			}
		})

	//! "Назад"
	av.v.pageRegister.buttonExit.SetSelectedFunc(func() {
		// switch
		av.v.pageMain.pages.SwitchToPage("login_page")
		// focus
		av.v.pageApp.app.SetInputCapture(av.v.pageLogin.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputEmail)
		// messageBox
		av.v.pageMain.messageBoxL.Clear()
		av.v.pageMain.messageBoxR.Clear()
		av.v.pageMain.statusBar.cellResponseStatus.SetText("")
	})

	//! form grid
	av.v.pageRegister.grid.
		SetRows(8, 1, 1).
		SetColumns(50).
		SetGap(1, 0).
		AddItem(av.v.pageRegister.form.form, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageRegister.buttonRegister, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageRegister.buttonExit, 2, 0, 1, 1, 0, 0, true)

	//! main grid
	av.v.pageRegister.mainGrid.
		SetRows(0, 20, 0).
		SetColumns(0, 40, 0).
		AddItem(av.v.pageRegister.grid, 1, 1, 1, 1, 0, 0, true)

	//! InputCapture
	av.v.pageRegister.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageRegister.form.inputEmail:
				av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputPassword)
			case av.v.pageRegister.form.inputPassword:
				av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputRPassword)
			case av.v.pageRegister.form.inputRPassword:
				av.v.pageApp.app.SetFocus(av.v.pageRegister.buttonRegister)
			case av.v.pageRegister.buttonRegister:
				av.v.pageApp.app.SetFocus(av.v.pageRegister.buttonExit)
			case av.v.pageRegister.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageRegister.form.inputEmail)
			}
			return nil
		}

		// if event.Key() == tcell.KeyEnter {
		// 	currentFocus := av.v.pageApp.App.GetFocus()
		// 	switch currentFocus {
		// 	case av.v.pageRegister.Form.Form.GetFormItem(2):
		// 		av.v.pageApp.App.SetFocus(av.v.pageRegister.ButtonRegister)
		// 	}
		// }

		return event
	}
}
