package appview

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataAddFile struct {
	treeview     *tview.TreeView
	buttonAdd    *tview.Button
	buttonCancel *tview.Button
	grid         *tview.Grid
	inputCapture func(event *tcell.EventKey) *tcell.EventKey
	page         *tview.Pages
}

func newDataAddFile() *pageDataAddFile {
	return &pageDataAddFile{
		treeview:     tview.NewTreeView(),
		buttonAdd:    tview.NewButton("Добавить"),
		buttonCancel: tview.NewButton("Отмена"),
		grid:         tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) VDataAddFile() {
	//! treeview
	var roots []string
	if runtime.GOOS == "windows" {
		for drive := 'A'; drive <= 'Z'; drive++ {
			drivePath := string(drive) + `:\`
			if _, err := os.Stat(drivePath); err == nil {
				roots = append(roots, drivePath)
			}
		}
	} else {
		roots = append(roots, "/")
	}

	root := tview.NewTreeNode("Выберите файл").
		SetColor(tview.Styles.TitleColor.TrueColor())
	av.v.pageData.pageDataAddFile.treeview.SetRoot(root).
		SetCurrentNode(root)

	for _, drive := range roots {
		driveNode := tview.NewTreeNode(drive).
			SetReference(drive).
			SetColor(tview.Styles.PrimaryTextColor)
		root.AddChild(driveNode)
	}

	add := func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			return
		}

		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name()))
			if file.IsDir() {
				node.SetColor(tview.Styles.SecondaryTextColor)
			} else {
				node.SetColor(tview.Styles.PrimaryTextColor)
			}
			target.AddChild(node)
		}
	}

	av.v.pageData.pageDataAddFile.treeview.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Выбор корневого узла ничего не делает
		}
		path := reference.(string)

		// Проверяем, является ли узел директорией
		info, err := os.Stat(path)
		if err != nil {
			return
		}

		if info.IsDir() {
			children := node.GetChildren()
			if len(children) == 0 {
				// Загружаем и отображаем содержимое директории
				add(node, path)
			} else {
				// Сворачиваем, если видно, разворачиваем, если свернуто
				node.SetExpanded(!node.IsExpanded())
			}
		} else {
			// Обработка выбора файла
		}
	})

	//! grid
	av.v.pageData.pageDataAddFile.grid.
		SetBorders(false).
		SetRows(0).
		SetColumns(0).
		AddItem(av.v.pageData.pageDataAddFile.treeview, 0, 0, 1, 1, 0, 0, true)
}
