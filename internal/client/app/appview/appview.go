package appview

import (
	"context"

	"github.com/golovanevvs/confidant/internal/client/model"
	"go.uber.org/zap"
)

type IServiceAccount interface {
	CreateAccount(ctx context.Context, email, password string) (registerAccountResp *model.AccountResp, err error)
	GetAccessToken(ctx context.Context, refreshTokenString string) (accessTokenString string, err error)
	Login(ctx context.Context, email, password string) (registerAccountResp *model.AccountResp, err error)
	LoginAtStart(ctx context.Context) (accountID int, email string, refreshTokenString string, err error)
	Logout(ctx context.Context) (err error)
}

type IServiceManage interface {
	GetServerStatus(ctx context.Context) (statusResp *model.StatusResp, err error)
}

type IServiceGroups interface {
	GetGroups(ctx context.Context, email string) (groups []model.Group, err error)
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
	GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error)
	AddEmail(ctx context.Context, groupID int, email string) (err error)
	GetGroupIDs(ctx context.Context, email string) (groupServerIDs map[int]struct{}, groupNoServerIDs map[int]struct{}, err error)
}

type IServiceData interface {
	GetDataTitles(ctx context.Context, accountID int, groupID int) (dataTitles []string, err error)
	GetDataTypes(ctx context.Context, accountID int, groupID int) (dataTypes []string, err error)
	GetDataIDAndType(ctx context.Context, groupID int, dataTitle string) (dataID int, dataType string, err error)
	AddNote(ctx context.Context, data model.NoteDec, accountID int, groupID int) (err error)
	GetNote(ctx context.Context, dataID int) (data *model.NoteDec, err error)
	AddPass(ctx context.Context, data model.PassDec, accountID int, groupID int) (err error)
	GetPass(ctx context.Context, dataID int) (data *model.PassDec, err error)
	AddCard(ctx context.Context, data model.CardDec, accountID int, groupID int) (err error)
	GetCard(ctx context.Context, dataID int) (data *model.CardDec, err error)
	AddFile(ctx context.Context, data model.FileDec, accountID int, groupID int, filepath string) (err error)
	GetFile(ctx context.Context, dataID int) (data *model.FileDec, err error)
	SaveToFile(ctx context.Context, dataID int, filepath string) (err error)
}

type IService interface {
	IServiceAccount
	IServiceManage
	IServiceGroups
	IServiceData
}

type view struct {
	pageApp      *pageApp
	pageMain     *pageMain
	pageLogin    *pageLogin
	pageRegister *pageRegister
	pageGroups   *pageGroups
	pageData     *pageData
}

type appView struct {
	sv           IService
	lg           *zap.SugaredLogger
	v            view
	accessToken  string
	refreshToken string
	account      model.Account
	groups       []model.Group
	groupID      int
	groupTitle   string
	dataTitles   []string
	dataTypes    []string
	dataTitle    string
	dataType     string
	dataID       int
	dataNoteID   int
	dataPassID   int
	dataCardID   int
	dataFileID   int
	dataFilepath string
	dataFilename string
}

func New(sv IService, lg *zap.SugaredLogger) *appView {
	return &appView{
		v: view{
			pageApp:      newPageApp(),
			pageMain:     newPageMain(),
			pageLogin:    newPageLogin(),
			pageRegister: newPageRegister(),
			pageGroups:   newPageGroups(),
			pageData:     newPageData(),
		},
		sv:           sv,
		lg:           lg,
		accessToken:  "",
		refreshToken: "",
		account: model.Account{
			ID:    -1,
			Email: "",
		},
		groups:       []model.Group{},
		groupID:      -1,
		groupTitle:   "",
		dataTitles:   []string{},
		dataTypes:    []string{},
		dataTitle:    "",
		dataType:     "",
		dataID:       -1,
		dataPassID:   -1,
		dataCardID:   -1,
		dataFileID:   -1,
		dataNoteID:   -1,
		dataFilepath: "",
		dataFilename: "",
	}
}

func (av *appView) Run() error {

	av.vMain()
	av.vLogin()
	av.vRegister()
	av.vGroups()
	av.vGroupsSelect()
	av.vGroupsAddGroup()
	av.vGroupsEditEmails()
	av.vData()
	av.vDataSelectType()
	av.vDataViewEmpty()
	av.vDataViewNote()
	av.vDataViewPass()
	av.vDataViewCard()
	av.vDataViewFile()
	av.vDataAddNote()
	av.vDataAddPass()
	av.vDataAddCard()
	av.vDataAddFile()
	av.vUpdateConnectionStatus()

	return av.vApp()
}
