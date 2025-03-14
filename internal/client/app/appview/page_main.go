package appview

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/rivo/tview"
)

type statusBar struct {
	cellTypeConnect     *tview.TableCell
	cellServerConnect   *tview.TableCell
	cellServerDBConnect *tview.TableCell
	cellResponseStatus  *tview.TableCell
	cellActiveAccount   *tview.TableCell
	table               *tview.Table
}

type pageMain struct {
	pages       *tview.Pages
	messageBoxL *tview.TextView
	messageBoxR *tview.TextView
	statusBar   statusBar
	mainGrid    *tview.Grid
}

func newPageMain() *pageMain {
	return &pageMain{
		pages:       tview.NewPages(),
		messageBoxL: tview.NewTextView(),
		messageBoxR: tview.NewTextView(),
		statusBar: statusBar{
			cellTypeConnect:     tview.NewTableCell(""),
			cellServerConnect:   tview.NewTableCell(""),
			cellServerDBConnect: tview.NewTableCell(""),
			cellResponseStatus:  tview.NewTableCell(""),
			cellActiveAccount:   tview.NewTableCell(""),
			table:               tview.NewTable(),
		},
		mainGrid: tview.NewGrid(),
	}
}

func (av *appView) vMain() {
	var statusServer *model.StatusResp
	var statusServerErr error

	//! left message box
	av.v.pageMain.messageBoxL.SetDynamicColors(true)
	av.v.pageMain.messageBoxL.SetTextAlign(tview.AlignLeft)
	av.v.pageMain.messageBoxL.SetBorder(true).SetBorderColor(tcell.ColorRed)
	av.v.pageMain.messageBoxL.SetTitle(" Сообщения ")

	//! right message box
	av.v.pageMain.messageBoxR.SetDynamicColors(true)
	av.v.pageMain.messageBoxR.SetTextAlign(tview.AlignLeft)
	av.v.pageMain.messageBoxR.SetBorder(true).SetBorderColor(tcell.ColorRed)
	av.v.pageMain.messageBoxR.SetTitle(" Дополнительная информация ")

	//! status bar
	av.v.pageMain.statusBar.table.SetBorders(true).SetBordersColor(tcell.ColorRed)
	av.v.pageMain.statusBar.table.SetCell(0, 0, tview.NewTableCell("Тип соединения").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.statusBar.table.SetCell(0, 1, tview.NewTableCell("Соединение с сервером").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.statusBar.table.SetCell(0, 2, tview.NewTableCell("Соединение с БД сервера").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.statusBar.table.SetCell(0, 3, tview.NewTableCell("Статус операции").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.statusBar.table.SetCell(0, 4, tview.NewTableCell("Активный аккаунт").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.statusBar.cellTypeConnect.SetText("[green]REST API").SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.statusBar.table.SetCell(1, 0, av.v.pageMain.statusBar.cellTypeConnect)
	av.v.pageMain.statusBar.cellServerConnect.SetText("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.statusBar.table.SetCell(1, 1, av.v.pageMain.statusBar.cellServerConnect)
	av.v.pageMain.statusBar.cellServerDBConnect.SetText("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.statusBar.table.SetCell(1, 2, av.v.pageMain.statusBar.cellServerDBConnect)
	av.v.pageMain.statusBar.cellResponseStatus.SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.statusBar.table.SetCell(1, 3, av.v.pageMain.statusBar.cellResponseStatus)
	av.v.pageMain.statusBar.cellActiveAccount.SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.statusBar.table.SetCell(1, 4, av.v.pageMain.statusBar.cellActiveAccount)

	//! updating status bar
	updateCellServerConnectChan := make(chan string)
	updateCellServerDBConnectChan := make(chan string)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			statusServer, statusServerErr = av.sv.GetServerStatus(context.Background())
			if statusServerErr != nil {
				updateCellServerConnectChan <- "[red]Отсутствует"
				updateCellServerDBConnectChan <- "[red]Отсутствует"
			} else {
				if statusServer.HTTPStatusCode == 200 {
					updateCellServerConnectChan <- "[green]OK"
					updateCellServerDBConnectChan <- "[green]OK"
				} else {
					updateCellServerConnectChan <- "[green]OK"
					updateCellServerDBConnectChan <- "[red]Отсутствует"
				}
			}
		}
	}()

	go func() {
		for info := range updateCellServerConnectChan {
			av.v.pageApp.app.QueueUpdateDraw(func() {
				av.v.pageMain.statusBar.cellServerConnect.SetText(info)
			})
		}
	}()
	go func() {
		for info := range updateCellServerDBConnectChan {
			av.v.pageApp.app.QueueUpdateDraw(func() {
				av.v.pageMain.statusBar.cellServerDBConnect.SetText(info)
			})
		}
	}()

	//! main grid
	av.v.pageMain.mainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	av.v.pageMain.mainGrid.SetRows(0, 8, 5)
	av.v.pageMain.mainGrid.AddItem(av.v.pageMain.pages, 0, 0, 1, 2, 0, 0, true)
	av.v.pageMain.mainGrid.AddItem(av.v.pageMain.messageBoxL, 1, 0, 1, 1, 0, 0, true)
	av.v.pageMain.mainGrid.AddItem(av.v.pageMain.messageBoxR, 1, 1, 1, 1, 0, 0, true)
	av.v.pageMain.mainGrid.AddItem(av.v.pageMain.statusBar.table, 2, 0, 1, 2, 0, 0, false)

	// av.v.pageMain.App.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	// 	focus := av.v.pageMain.App.GetFocus()
	// 	if focus != av.v.pageGroups.ButtonDeleteEmail {
	// 		return event, action

	// 	}
	// 	return nil, 0
	// })

	//! adding pages
	av.v.pageMain.pages.AddPage("groups_page", av.v.pageGroups.gridMain, true, true)
	av.v.pageMain.pages.AddPage("register_page", av.v.pageRegister.mainGrid, true, true)
	av.v.pageMain.pages.AddAndSwitchToPage("login_page", av.v.pageLogin.mainGrid, true)

}
