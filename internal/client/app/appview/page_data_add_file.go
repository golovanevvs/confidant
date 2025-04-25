package appview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/rivo/tview"
)

type pageDataAddFile struct {
	treeview          *tview.TreeView
	textviewFileNameL *tview.TextView
	textviewFileSizeL *tview.TextView
	textviewFileDateL *tview.TextView
	textviewDescL     *tview.TextView
	textviewTitleL    *tview.TextView
	textviewFileName  *tview.TextView
	textviewFileSize  *tview.TextView
	textviewFileDate  *tview.TextView
	textareaDesc      *tview.TextArea
	textareaTitle     *tview.TextArea
	buttonAdd         *tview.Button
	buttonCancel      *tview.Button
	gridData          *tview.Grid
	gridButtons       *tview.Grid
	grid              *tview.Grid
	inputCapture      func(event *tcell.EventKey) *tcell.EventKey
	page              *tview.Pages
}

func newDataAddFile() *pageDataAddFile {
	return &pageDataAddFile{
		treeview:          tview.NewTreeView(),
		textviewFileNameL: tview.NewTextView(),
		textviewFileSizeL: tview.NewTextView(),
		textviewFileDateL: tview.NewTextView(),
		textviewDescL:     tview.NewTextView(),
		textviewTitleL:    tview.NewTextView(),
		textviewFileName:  tview.NewTextView(),
		textviewFileSize:  tview.NewTextView(),
		textviewFileDate:  tview.NewTextView(),
		textareaDesc:      tview.NewTextArea(),
		textareaTitle:     tview.NewTextArea(),
		buttonAdd:         tview.NewButton("Добавить"),
		buttonCancel:      tview.NewButton("Отмена"),
		gridData:          tview.NewGrid(),
		gridButtons:       tview.NewGrid(),
		grid:              tview.NewGrid(),
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		page: tview.NewPages(),
	}
}

func (av *appView) vDataAddFile() {
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
			return // nothing happens when selecting a mountain node
		}
		path := reference.(string)

		// checking whether a node is a directory
		info, err := os.Stat(path)
		if err != nil {
			return
		}

		if info.IsDir() {
			children := node.GetChildren()
			if len(children) == 0 {
				// displaying the contents of a directory
				add(node, path)
			} else {
				// folding, unfolding
				node.SetExpanded(!node.IsExpanded())
			}
		} else {
			// file selection processing
			av.dataFilename = info.Name()
			av.dataFilepath = path
			av.v.pageData.pageDataAddFile.textviewFileName.SetText(info.Name())
			av.v.pageData.pageDataAddFile.textviewFileSize.SetText(fmt.Sprintf("%d байт", info.Size()))
			av.v.pageData.pageDataAddFile.textviewFileDate.SetText(info.ModTime().Format("02.01.2006 15:04:05"))
		}
	})

	// processing changes to a dedicated node
	av.v.pageData.pageDataAddFile.treeview.SetChangedFunc(func(node *tview.TreeNode) {
		av.dataFilename = ""
		av.dataFilepath = ""
		av.v.pageData.pageDataAddFile.textviewFileName.SetText("")
		av.v.pageData.pageDataAddFile.textviewFileSize.SetText("")
		av.v.pageData.pageDataAddFile.textviewFileDate.SetText("")
	})

	//! label names
	av.v.pageData.pageDataAddFile.textviewFileNameL.SetText("Выбран файл:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddFile.textviewFileSizeL.SetText("Размер:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddFile.textviewFileDateL.SetText("Дата изменения:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddFile.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataAddFile.textviewTitleL.SetText("Название:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! Добавить
	av.v.pageData.pageDataAddFile.buttonAdd.SetSelectedFunc(func() {
		title := av.v.pageData.pageDataAddFile.textareaTitle.GetText()
		desc := av.v.pageData.pageDataAddFile.textareaDesc.GetText()
		filesize := av.v.pageData.pageDataAddFile.textviewFileSize.GetText(false)
		filedate := av.v.pageData.pageDataAddFile.textviewFileDate.GetText(false)
		data := model.FileDec{
			Title:    title,
			Desc:     desc,
			Filename: av.dataFilename,
			Filesize: filesize,
			Filedate: filedate,
		}
		err := av.sv.AddFile(context.Background(), data, av.account.ID, av.groupID, av.dataFilepath)
		if err != nil {
			av.v.pageMain.messageBoxL.SetText(err.Error())
		} else {
			av.aPageDataSwitch()
		}
	})

	//! Отмена
	av.v.pageData.pageDataAddFile.buttonCancel.SetSelectedFunc(func() {
		av.aPageDataSwitch()
	})

	//! data grid
	av.v.pageData.pageDataAddFile.gridData.
		SetBorders(true).
		SetRows(0, 1, 1, 1, 1).
		SetColumns(15, 15, 15, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataAddFile.treeview, 0, 0, 1, 4, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewFileNameL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewFileName, 1, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewFileSizeL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewFileSize, 2, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewFileDateL, 2, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewFileDate, 2, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewDescL, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textareaDesc, 3, 1, 1, 3, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textviewTitleL, 4, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.textareaTitle, 4, 1, 1, 3, 0, 0, true)

	//! buttons grid
	av.v.pageData.pageDataAddFile.gridButtons.
		SetBorders(false).
		SetRows(1).
		SetColumns(10, 10).
		SetGap(0, 1).
		AddItem(av.v.pageData.pageDataAddFile.buttonAdd, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.buttonCancel, 0, 1, 1, 1, 0, 0, true)

	//! grid
	av.v.pageData.pageDataAddFile.grid.
		SetBorders(false).
		SetRows(0, 1).
		SetColumns(0).
		AddItem(av.v.pageData.pageDataAddFile.gridData, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataAddFile.gridButtons, 1, 0, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataAddFile.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.pageDataAddFile.treeview:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddFile.textareaDesc)
			case av.v.pageData.pageDataAddFile.textareaDesc:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddFile.textareaTitle)
			case av.v.pageData.pageDataAddFile.textareaTitle:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddFile.buttonAdd)
			case av.v.pageData.pageDataAddFile.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddFile.buttonCancel)
			case av.v.pageData.pageDataAddFile.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataAddFile.treeview)
			}
			return nil
		}
		return event
	}
}
