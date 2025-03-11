package appview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/rivo/tview"
	"go.uber.org/zap"
)

type IAccountService interface {
	RegisterAccount(email, password string) (registerAccountResp *model.RegisterAccountResp, err error)
	// Login(email, password string) (accountID int, err error)
	// ChangePassword(email, password, newPassword string) error
	//GetUser(login string) (string, error)
}

type IStatusServerService interface {
	GetServerStatus() (statusResp *model.StatusResp, err error)
}

type IService interface {
	IAccountService
	IStatusServerService
}

type View struct {
	pageMain     PageMain
	pageLogin    PageLogin
	pageRegister PageRegister
	pageGroups   PageGroups
}

type AppView struct {
	sv IService
	lg *zap.SugaredLogger
	v  View
}

func New(sv IService, lg *zap.SugaredLogger) *AppView {
	return &AppView{
		v: View{
			pageMain: PageMain{
				App:         tview.NewApplication(),
				Pages:       tview.NewPages(),
				MessageBoxL: tview.NewTextView(),
				MessageBoxR: tview.NewTextView(),
				StatusBar: StatusBar{
					CellTypeConnect:     tview.NewTableCell(""),
					CellServerConnect:   tview.NewTableCell(""),
					CellServerDBConnect: tview.NewTableCell(""),
					CellResponseStatus:  tview.NewTableCell(""),
					CellActiveAccount:   tview.NewTableCell(""),
					Table:               tview.NewTable(),
				},
				MainGrid: tview.NewGrid(),
			},
			pageLogin: PageLogin{
				InputCapture: func(event *tcell.EventKey) *tcell.EventKey {
					return event
				},
				Form: FormPageLogin{
					Form:          tview.NewForm(),
					InputEmail:    tview.NewInputField(),
					InputPassword: tview.NewInputField(),
				},
				ButtonLogin:    tview.NewButton("Войти"),
				ButtonRegister: tview.NewButton("Регистрация"),
				ButtonExit:     tview.NewButton("Выход"),
				Grid:           tview.NewGrid(),
				MainGrid:       tview.NewGrid(),
			},
			pageRegister: PageRegister{
				InputCapture: func(event *tcell.EventKey) *tcell.EventKey {
					return event
				},
				Form: FormPageRegister{
					Form:           tview.NewForm(),
					InputEmail:     tview.NewInputField(),
					InputPassword:  tview.NewInputField(),
					InputRPassword: tview.NewInputField(),
				},
				ButtonRegister: tview.NewButton("Зарегистрироваться"),
				ButtonExit:     tview.NewButton("Назад"),
				Grid:           tview.NewGrid(),
				MainGrid:       tview.NewGrid(),
			},
			pageGroups: PageGroups{
				ListGroups: tview.NewList(),
				ListEmails: tview.NewList(),
				GridMain:   tview.NewGrid(),
				PageSelect: PageSelect{
					GridSelectButtons: tview.NewGrid(),
					ButtonSelect:      tview.NewButton("Выбрать группу"),
					ButtonNew:         tview.NewButton("Создать группу"),
					ButtonSettings:    tview.NewButton("Настроить группу"),
					ButtonDelete:      tview.NewButton("Удалить группу"),
					ButtonLogout:      tview.NewButton("Выйти из аккаунта"),
					ButtonExit:        tview.NewButton("Выход"),
					Grid:              tview.NewGrid(),
					InputCapture: func(event *tcell.EventKey) *tcell.EventKey {
						return event
					},
					Page: tview.NewPages(),
				},
				PageEdit: PageEdit{
					FormAddEmail: FormAdd{
						InputEmail: tview.NewInputField(),
						Form:       tview.NewForm(),
					},
					ButtonAddEmail:    tview.NewButton("Добавить e-mail"),
					ButtonDeleteEmail: tview.NewButton("Удалить e-mail"),
					ButtonEditExit:    tview.NewButton("Назад"),
					Grid:              tview.NewGrid(),
					InputCapture: func(event *tcell.EventKey) *tcell.EventKey {
						return event
					},
					Page: tview.NewPages(),
				},
				PagesSelEd: tview.NewPages(),
			},
		},
		sv: sv,
		lg: lg,
	}
}
