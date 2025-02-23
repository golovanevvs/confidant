package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/app/appview"
	"github.com/rivo/tview"
)

func RunApp() {
	//! Начало
	av := appview.AppView{}

	// Приложение
	app := tview.NewApplication()

	// Контейнер страниц
	pages := tview.NewPages()

	mainGrid := tview.NewGrid()
	mainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	mainGrid.AddItem(pages, 0, 0, 1, 1, 0, 0, true)

	//! Добавление страницы

	pages.AddAndSwitchToPage("login_page", av.LoginPage(app, pages), true)

	//! Запуск приложения
	app.SetRoot(mainGrid, true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
