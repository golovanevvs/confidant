package appview

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type pageDataViewFile struct {
	treeview          *tview.TreeView
	textviewFileNameL *tview.TextView
	textviewFilePathL *tview.TextView
	textviewFileSizeL *tview.TextView
	textviewFileDateL *tview.TextView
	textviewDescL     *tview.TextView
	textviewFileName  *tview.TextView
	textviewFilePath  *tview.TextView
	textviewFileSize  *tview.TextView
	textviewFileDate  *tview.TextView
	textviewDesc      *tview.TextView
	buttonSave        *tview.Button
	buttonCancel      *tview.Button
	gridData          *tview.Grid
	gridButtons       *tview.Grid
	grid              *tview.Grid
	inputCapture      func(event *tcell.EventKey) *tcell.EventKey
	page              *tview.Pages
}

func newDataViewFile() *pageDataViewFile {
	return &pageDataViewFile{
		treeview:          tview.NewTreeView(),
		textviewFileNameL: tview.NewTextView(),
		textviewFilePathL: tview.NewTextView(),
		textviewFileSizeL: tview.NewTextView(),
		textviewFileDateL: tview.NewTextView(),
		textviewDescL:     tview.NewTextView(),
		textviewFileName:  tview.NewTextView(),
		textviewFilePath:  tview.NewTextView(),
		textviewFileSize:  tview.NewTextView(),
		textviewFileDate:  tview.NewTextView(),
		textviewDesc:      tview.NewTextView(),
		buttonSave:        tview.NewButton("Сохранить"),
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

func (av *appView) vDataViewFile() {
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

	root := tview.NewTreeNode("Выберите путь для сохранения файла").
		SetColor(tview.Styles.TitleColor.TrueColor())
	av.v.pageData.pageDataViewFile.treeview.SetRoot(root).
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
			if file.IsDir() {
				node := tview.NewTreeNode(file.Name()).
					SetReference(filepath.Join(path, file.Name())).
					SetColor(tview.Styles.SecondaryTextColor)
				target.AddChild(node)
			}
		}
	}

	av.v.pageData.pageDataViewFile.treeview.SetSelectedFunc(func(node *tview.TreeNode) {
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
			av.v.pageData.pageDataViewFile.textviewFilePath.SetText(path)
		} else {
			// file selection processing
			// av.v.pageData.pageDataViewFile.textviewFileName.SetText(info.Name())
			// av.v.pageData.pageDataViewFile.textviewFileSize.SetText(fmt.Sprintf("%d байт", info.Size()))
			// av.v.pageData.pageDataViewFile.textviewFileDate.SetText(info.ModTime().Format("02.01.2006 15:04:05"))
		}
	})

	// processing changes to a dedicated node
	av.v.pageData.pageDataViewFile.treeview.SetChangedFunc(func(node *tview.TreeNode) {
		av.v.pageData.pageDataViewFile.textviewFilePath.SetText("")
		av.v.pageData.pageDataViewFile.textviewFileSize.SetText("")
		av.v.pageData.pageDataViewFile.textviewFileDate.SetText("")
	})

	//! label names
	av.v.pageData.pageDataViewFile.textviewFileNameL.SetText("Имя файла:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewFile.textviewFilePathL.SetText("Выбран путь:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewFile.textviewFileSizeL.SetText("Размер:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewFile.textviewFileDateL.SetText("Дата изменения:").
		SetTextColor(av.v.pageApp.colorTitle)
	av.v.pageData.pageDataViewFile.textviewDescL.SetText("Описание:").
		SetTextColor(av.v.pageApp.colorTitle)

	//! Сохранить
	av.v.pageData.pageDataViewFile.buttonSave.SetSelectedFunc(func() {
		path := av.v.pageData.pageDataViewFile.textviewFilePath.GetText(false)
		filename := av.v.pageData.pageDataViewFile.textviewFileName.GetText(false)
		fp := filepath.Join(path, filename)
		err := av.sv.SaveToFile(context.Background(), av.dataID, fp)
		if err != nil {
			av.v.pageMain.messageBoxL.SetText(err.Error())
		} else {
			av.aPageDataSwitch()
		}
	})

	//! Отмена
	av.v.pageData.pageDataViewFile.buttonCancel.SetSelectedFunc(func() {
		av.aPageDataSwitch()
	})

	//! data grid
	av.v.pageData.pageDataViewFile.gridData.
		SetBorders(true).
		SetRows(0, 1, 1, 1).
		SetColumns(15, 20, 15, 0).
		SetGap(1, 0).
		AddItem(av.v.pageData.pageDataViewFile.treeview, 0, 0, 1, 4, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFileNameL, 1, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFileName, 1, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFilePathL, 1, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFilePath, 1, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFileSizeL, 2, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFileSize, 2, 1, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFileDateL, 2, 2, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewFileDate, 2, 3, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewDescL, 3, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.textviewDesc, 3, 1, 1, 3, 0, 0, true)

	//! buttons grid
	av.v.pageData.pageDataViewFile.gridButtons.
		SetBorders(false).
		SetRows(1).
		SetColumns(10, 10).
		SetGap(0, 1).
		AddItem(av.v.pageData.pageDataViewFile.buttonSave, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.buttonCancel, 0, 1, 1, 1, 0, 0, true)

	//! grid
	av.v.pageData.pageDataViewFile.grid.
		SetBorders(false).
		SetRows(0, 1).
		SetColumns(0).
		AddItem(av.v.pageData.pageDataViewFile.gridData, 0, 0, 1, 1, 0, 0, true).
		AddItem(av.v.pageData.pageDataViewFile.gridButtons, 1, 0, 1, 1, 0, 0, true)

	//! inputCapture
	av.v.pageData.pageDataViewFile.inputCapture = func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			currentFocus := av.v.pageApp.app.GetFocus()
			switch currentFocus {
			case av.v.pageData.listTitles:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonAdd)
			case av.v.pageData.buttonAdd:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonEdit)
			case av.v.pageData.buttonEdit:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonDelete)
			case av.v.pageData.buttonDelete:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonBack)
			case av.v.pageData.buttonBack:
				av.v.pageApp.app.SetFocus(av.v.pageData.buttonExit)
			case av.v.pageData.buttonExit:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewFile.treeview)
			case av.v.pageData.pageDataViewFile.treeview:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewFile.buttonSave)
			case av.v.pageData.pageDataViewFile.buttonSave:
				av.v.pageApp.app.SetFocus(av.v.pageData.pageDataViewFile.buttonCancel)
			case av.v.pageData.pageDataViewFile.buttonCancel:
				av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
			}
			return nil
		}
		return event
	}
}

func (av *appView) vPageDataViewFileUpdate() {
	av.vClearMessages()

	data, err := av.sv.GetFile(context.Background(), av.dataID)
	if err != nil {
		av.v.pageMain.messageBoxL.SetText(err.Error())
	} else {
		av.v.pageData.pageDataViewFile.textviewFileName.SetText(data.Filename)
		av.v.pageData.pageDataViewFile.textviewDesc.SetText(data.Desc)
		av.v.pageData.pageDataViewFile.textviewFileSize.SetText(data.Filesize)
		av.v.pageData.pageDataViewFile.textviewFileDate.SetText(data.Filedate)
		av.v.pageMain.pages.SwitchToPage("data_page")
		av.v.pageData.pages.SwitchToPage("data_view_file_page")
		av.v.pageApp.app.SetInputCapture(av.v.pageData.pageDataViewFile.inputCapture)
		av.v.pageApp.app.SetFocus(av.v.pageData.listTitles)
	}
}
