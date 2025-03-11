package appview

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/rivo/tview"
)

func (av *AppView) Run() error {
	//!Beginning
	action := "run"

	// app
	av.v.pageMain.App = tview.NewApplication()

	// page container
	av.v.pageMain.Pages = tview.NewPages()

	av.VMain()

	//! LOGIN PAGE

	av.VLogin()

	//! REGISTER PAGE

	av.VRegister()

	//! GROUPS PAGE

	av.VGroups()

	//! Launching the app
	av.v.pageMain.App.SetRoot(av.v.pageMain.MainGrid, true)

	av.v.pageMain.App.SetInputCapture(av.v.pageLogin.InputCapture)
	av.v.pageMain.App.SetFocus(av.v.pageLogin.Form.InputEmail)

	if err := av.v.pageMain.App.Run(); err != nil {
		return fmt.Errorf("%s: %s: %w: %w", customerrors.ClientAppViewErr, action, customerrors.ErrRunAppView, err)
	}

	return nil
}
