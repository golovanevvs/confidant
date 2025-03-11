package appview

import "github.com/rivo/tview"

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
}
