package appview

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/rivo/tview"
)

type StatusBar struct {
	CellTypeConnect     *tview.TableCell
	CellServerConnect   *tview.TableCell
	CellServerDBConnect *tview.TableCell
	CellResponseStatus  *tview.TableCell
	CellActiveAccount   *tview.TableCell
	Table               *tview.Table
}

type PageMain struct {
	App         *tview.Application
	Pages       *tview.Pages
	MessageBoxL *tview.TextView
	MessageBoxR *tview.TextView
	StatusBar   StatusBar
	MainGrid    *tview.Grid
}

func (av *AppView) VMain() {
	// action := "run"

	var statusServer *model.StatusResp
	var statusServerErr error

	//? left message box
	av.v.pageMain.MessageBoxL.SetDynamicColors(true)
	av.v.pageMain.MessageBoxL.SetTextAlign(tview.AlignLeft)
	av.v.pageMain.MessageBoxL.SetBorder(true).SetBorderColor(tcell.ColorRed)
	av.v.pageMain.MessageBoxL.SetTitle(" Сообщения ")

	//? right message box
	av.v.pageMain.MessageBoxR.SetDynamicColors(true)
	av.v.pageMain.MessageBoxR.SetTextAlign(tview.AlignLeft)
	av.v.pageMain.MessageBoxR.SetBorder(true).SetBorderColor(tcell.ColorRed)
	av.v.pageMain.MessageBoxR.SetTitle(" Дополнительная информация ")

	//? status bar
	av.v.pageMain.StatusBar.Table.SetBorders(true).SetBordersColor(tcell.ColorRed)
	av.v.pageMain.StatusBar.Table.SetCell(0, 0, tview.NewTableCell("Тип соединения").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.StatusBar.Table.SetCell(0, 1, tview.NewTableCell("Соединение с сервером").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.StatusBar.Table.SetCell(0, 2, tview.NewTableCell("Соединение с БД сервера").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.StatusBar.Table.SetCell(0, 3, tview.NewTableCell("Статус операции").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.StatusBar.Table.SetCell(0, 4, tview.NewTableCell("Активный аккаунт").SetAlign(tview.AlignCenter).SetExpansion(1))
	av.v.pageMain.StatusBar.CellTypeConnect.SetText("[green]REST API").SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.StatusBar.Table.SetCell(1, 0, av.v.pageMain.StatusBar.CellTypeConnect)
	av.v.pageMain.StatusBar.CellServerConnect.SetText("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.StatusBar.Table.SetCell(1, 1, av.v.pageMain.StatusBar.CellServerConnect)
	av.v.pageMain.StatusBar.CellServerDBConnect.SetText("[red]Отсутствует").SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.StatusBar.Table.SetCell(1, 2, av.v.pageMain.StatusBar.CellServerDBConnect)
	av.v.pageMain.StatusBar.CellResponseStatus.SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.StatusBar.Table.SetCell(1, 3, av.v.pageMain.StatusBar.CellResponseStatus)
	av.v.pageMain.StatusBar.CellActiveAccount.SetAlign(tview.AlignCenter).SetExpansion(1)
	av.v.pageMain.StatusBar.Table.SetCell(1, 4, av.v.pageMain.StatusBar.CellActiveAccount)

	updateCellServerConnectChan := make(chan string)
	updateCellServerDBConnectChan := make(chan string)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			statusServer, statusServerErr = av.sv.GetServerStatus()
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
			av.v.pageMain.App.QueueUpdateDraw(func() {
				av.v.pageMain.StatusBar.CellServerConnect.SetText(info)
			})
		}
	}()
	go func() {
		for info := range updateCellServerDBConnectChan {
			av.v.pageMain.App.QueueUpdateDraw(func() {
				av.v.pageMain.StatusBar.CellServerDBConnect.SetText(info)
			})
		}
	}()

	//? main grid
	av.v.pageMain.MainGrid.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" Клиент [blue]системы [red]confidant ")
	av.v.pageMain.MainGrid.SetRows(0, 8, 5)
	av.v.pageMain.MainGrid.AddItem(av.v.pageMain.Pages, 0, 0, 1, 2, 0, 0, true)
	av.v.pageMain.MainGrid.AddItem(av.v.pageMain.MessageBoxL, 1, 0, 1, 1, 0, 0, true)
	av.v.pageMain.MainGrid.AddItem(av.v.pageMain.MessageBoxR, 1, 1, 1, 1, 0, 0, true)
	av.v.pageMain.MainGrid.AddItem(av.v.pageMain.StatusBar.Table, 2, 0, 1, 2, 0, 0, false)

	// av.v.pageMain.App.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
	// 	focus := av.v.pageMain.App.GetFocus()
	// 	if focus != av.v.pageGroups.ButtonDeleteEmail {
	// 		return event, action

	// 	}
	// 	return nil, 0
	// })

	//! Adding Pages

	av.v.pageMain.Pages.AddPage("groups_page", av.v.pageGroups.GridMain, true, true)
	av.v.pageMain.Pages.AddPage("register_page", av.v.pageRegister.MainGrid, true, true)
	av.v.pageMain.Pages.AddAndSwitchToPage("login_page", av.v.pageLogin.MainGrid, true)
}
