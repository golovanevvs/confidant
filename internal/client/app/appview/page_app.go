package appview

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

type pageApp struct {
	app *tview.Application
}

func newPageApp() *pageApp {
	return &pageApp{
		app: tview.NewApplication(),
	}
}

func (av *appView) vApp() error {
	action := "run"

	av.v.pageApp.app.SetRoot(av.v.pageMain.mainGrid, true)

	av.v.pageApp.app.SetInputCapture(av.v.pageLogin.inputCapture)
	av.v.pageApp.app.SetFocus(av.v.pageLogin.form.inputEmail)

	if err := av.v.pageApp.app.Run(); err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientAppViewErr,
			action,
			customerrors.ErrRunAppView,
			err)
	}
	return nil
}
