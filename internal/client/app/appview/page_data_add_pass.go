package appview

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/rivo/tview"
)

type pageDataAddPass struct {
	textviewLoginL *tview.TextView
	textviewPassL  *tview.TextView
	textviewDescL  *tview.TextView
	textviewTitleL *tview.TextView
	textareaLogin  *tview.TextArea
	textareaPass   *tview.TextArea
	textareaDesc   *tview.TextArea
	textareaTitle  *tview.TextArea
	buttonAdd      *tview.Button
	buttonCancel   *tview.Button
	gridData       *tview.Grid
	gridButtons    *tview.Grid
	grid           *tview.Grid
	inputCapture   func(event *tcell.EventKey) *tcell.EventKey
	page           *tview.Pages
}

func newPageDataAddPass() *pageDataAddPass {
	return &pageDataAddPass{
		textviewLoginL: tview.NewTextView(),
		textviewPassL:  tview.NewTextView(),
		textviewDescL:  tview.NewTextView(),
		textviewTitleL: tview.NewTextView(),
		textareaLogin:  tview.NewTextArea(),
		textareaPass:   tview.NewTextArea(),
		textareaDesc:   tview.NewTextArea(),
		textareaTitle:  tview.NewTextArea(),
		buttonAdd:      tview.NewButton("Добавить"),
		buttonCancel:   tview.NewButton("Отмена"),
		gridData:       tview.NewGrid(),
		gridButtons:    tview.NewGrid(),
		grid:           tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataAddPass() {
	//! label names
	av.v.pageData.pageDataAddPass.textviewLoginL.SetText("Логин:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddPass.textviewPassL.SetText("Пароль:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddPass.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddPass.textviewTitleL.SetText("Название:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! data grid
	av.v.pageData.pageDataAddPass.gridData.
		SetBorders(true).
		SetRows(1, 1, 0, 1).
		SetColumns(9, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataAddPass.textviewLoginL, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textviewPassL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textviewDescL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaLogin, 0, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaPass, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaDesc, 2, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textviewTitleL, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.textareaTitle, 3, 1, 1, 1, 0, 0, true)

	//! Добавить
	av.v.pageData.pageDataAddPass.buttonAdd.SetSelectedFunc(func() {
		title := av.v.pageData.pageDataAddPass.textareaTitle.GetText()
		desc := av.v.pageData.pageDataAddPass.textareaDesc.GetText()
		login := av.v.pageData.pageDataAddPass.textareaLogin.GetText()
		pass := av.v.pageData.pageDataAddPass.textareaPass.GetText()
		data := model.PassDec{
			Title: title,
			Desc:  desc,
			Login: login,
			Pass:  pass,
		}
		err := av.sv.AddPass(context.Background(), data, av.account.ID, av.groupID)
		if err != nil {
			av.v.pageMain.messageBoxL.SetText(err.Error())
		} else {
			av.aPageDataSwitch()
		}
	})

	//! Отмена
	av.v.pageData.pageDataAddPass.buttonCancel.SetSelectedFunc(func() {
		av.aPageDataSwitch()
	})

	//! buttons grid
	av.v.pageData.pageDataAddPass.gridButtons.
		SetBorders(false).
		SetRows(1).
		SetColumns(10, 10).
		SetGap(0, 1).
		AddItem(av.v.pageData.pageDataAddPass.buttonAdd, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.buttonCancel, 0, 1, 1, 1, 0, 0, true)

	//! grid
	av.v.pageData.pageDataAddPass.grid.
		SetBorders(false).
		SetRows(0, 1).
		SetColumns(0).
		AddItem(av.v.pageData.pageDataAddPass.gridData, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddPass.gridButtons, 1, 0, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataAddPass.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.pageDataAddPass.textareaLogin:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.textareaPass)
			case av.v.pageData.pageDataAddPass.textareaPass:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.textareaDesc)
			case av.v.pageData.pageDataAddPass.textareaDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.textareaTitle)
			case av.v.pageData.pageDataAddPass.textareaTitle:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.buttonAdd)
			case av.v.pageData.pageDataAddPass.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.buttonCancel)
			case av.v.pageData.pageDataAddPass.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddPass.textareaLogin)
			}
			return nil
		}
		return event
	}

}
