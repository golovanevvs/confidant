package appview

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	trhttp "github.com/golovanevvs/confidant/internal/client/transport/http"
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

	//message box
	messageBox := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("Сообщения, ошибки и прочее").
		SetBorder(true).
		SetBorderColor(tcell.ColorRed).
		SetTitle(" Сообщения ")

	mainGrid := tview.NewGrid()
	mainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	mainGrid.SetRows(0, 6)
	mainGrid.AddItem(pages, 0, 0, 1, 1, 0, 0, true)
	mainGrid.AddItem(messageBox, 1, 0, 1, 1, 0, 0, true)

	//! LOGIN PAGE
	formLoginPage := tview.NewForm()
	formLoginPage.SetHorizontal(false)
	formLoginPage.AddInputField("E-mail:", "", 0, nil, nil)
	formLoginPage.AddPasswordField("Пароль:", "", 0, '*', nil)
	buttonLoginLoginPage := tview.NewButton("Войти")
	buttonRegisterLoginPage := tview.NewButton("Зарегистрироваться").SetSelectedFunc(func() {
		pages.SwitchToPage("register_page")
	})
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

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

		if event.Key() == tcell.KeyEnter {
			currentFocus := app.GetFocus()
			switch currentFocus {
			case formLoginPage.GetFormItem(1):
				app.SetFocus(buttonLoginLoginPage)
			}
		}

		return event
	})

	//! REGISTER PAGE
	formRegisterPage := tview.NewForm()
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
				modal := tview.NewModal().
					SetText("Пароли не совпадают! Повторите ввод").
					AddButtons([]string{"Закрыть"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonLabel == "Закрыть" {
							pages.SwitchToPage("register_page")
						}
					})
				pages.AddPage("register_page_modal", modal, true, true)
			} else {
				account := model.Account{
					Email:    formRegisterPage.GetFormItem(0).(*tview.InputField).GetText(),
					Password: formRegisterPage.GetFormItem(1).(*tview.InputField).GetText(),
				}
				tr := trhttp.New()
				tr.RegisterAccount(account)
			}
		})
	buttonExitRegisterPage := tview.NewButton("Назад").SetSelectedFunc(func() {
		pages.AddAndSwitchToPage("login_page", mainGridLoginPage, true)
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

	// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Key() == tcell.KeyTAB {
	// 		currentFocus := app.GetFocus()
	// 		switch currentFocus {
	// 		case formRegisterPage.GetFormItem(0):
	// 			app.SetFocus(formRegisterPage.GetFormItem(1))
	// 		case formRegisterPage.GetFormItem(1):
	// 			app.SetFocus(formRegisterPage.GetFormItem(2))
	// 		case formRegisterPage.GetFormItem(2):
	// 			app.SetFocus(buttonRegisterRegisterPage)
	// 		case buttonRegisterRegisterPage:
	// 			app.SetFocus(buttonExitRegisterPage)
	// 		case buttonExitRegisterPage:
	// 			app.SetFocus(formRegisterPage.GetFormItem(0))
	// 		}
	// 		return nil
	// 	}

	// 	if event.Key() == tcell.KeyEnter {
	// 		currentFocus := app.GetFocus()
	// 		switch currentFocus {
	// 		case formRegisterPage.GetFormItem(2):
	// 			app.SetFocus(buttonRegisterRegisterPage)
	// 		}
	// 	}

	// 	return event
	// })

	//! Adding Pages

	pages.AddPage("register_page", mainGridRegisterPage, true, true)
	pages.AddAndSwitchToPage("login_page", mainGridLoginPage, true)

	//! Launching the app
	app.SetRoot(mainGrid, true)

	if err := app.Run(); err != nil {
		return fmt.Errorf("%s: %s: %w: %w", customerrors.AppViewErr, action, customerrors.ErrRunAppView, err)
	}

	return nil
}
