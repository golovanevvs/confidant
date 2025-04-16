package service_sync

import (
	"context"
	"time"

	"github.com/golovanevvs/confidant/internal/client/model"
)

type IServiceGroups interface {
	GetGroups(ctx context.Context, email string) (groups []model.Group, err error)
	GetGroupsByIDs(ctx context.Context, groupIDs []int) (groups []model.Group, err error)
	AddGroup(ctx context.Context, account *model.Account, title string) (err error)
	GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error)
	AddEmail(ctx context.Context, groupID int, email string) (err error)
	GetGroupIDs(ctx context.Context, email string) (groupServerIDs []int, groupNoServerIDs []int, err error)
	AddGroupBySync(ctx context.Context, group model.Group) (err error)
	UpdateGroupIDsOnServer(ctx context.Context, newGroupIDs map[int]int) (err error)
	GetEmails(ctx context.Context, groupIDs []int) (mapGroupIDEmails map[int][]string, err error)
	AddEmailsBySync(ctx context.Context, mapGroupIDEmails map[int][]string) (err error)
}

type IServiceData interface {
	GetDataIDs(ctx context.Context, groupIDs []int) (dataServerIDs []int, dataNoServerIDs []int, err error)
	GetDataDates(ctx context.Context, dataIDs []int) (dataDatesFromClient map[int]time.Time, err error)
	SaveDatas(ctx context.Context, datas []model.Data) (err error)
	SaveDataFile(ctx context.Context, dataID int, file []byte) (err error)
	GetDatas(ctx context.Context, dataIDs []int) (datas []model.Data, err error)
	UpdateDataIDsOnServer(ctx context.Context, newDataIDs map[int]int) (err error)
	GetDataFile(ctx context.Context, dataID int) (idOnServer int, file []byte, err error)
}

type ITransportSync interface {
	GetGroupIDs(ctx context.Context, accessToken string) (trResponse *model.GroupSyncResp, err error)
	GetGroups(ctx context.Context, accessToken string, groupIDs []int) (groupsFromServer []model.Group, err error)
	SendGroups(ctx context.Context, accessToken string, groups []model.Group) (groupIDs map[int]int, err error)
	GetDataIDs(ctx context.Context, accessToken string) (trResponse *model.DataSyncResp, err error)
	GetDataDates(ctx context.Context, accessToken string, dataIDs []int) (dataDatesFromServer map[int]time.Time, err error)
	GetDatas(ctx context.Context, accessToken string, dataIDs []int) (datasFromServer []model.Data, err error)
	GetDataFile(ctx context.Context, accessToken string, dataID int) (fileFromServer []byte, err error)
	SendDatas(ctx context.Context, accessToken string, datas []model.Data) (dataIDs map[int]int, err error)
	SendFile(ctx context.Context, accessToken string, dataID int, file []byte) (err error)
	GetEmails(ctx context.Context, accessToken string, groupIDs []int) (mapGroupIDEmailsFromServer map[int][]string, err error)
	SendEmails(ctx context.Context, accessToken string, mapGroupIDEmails map[int][]string) (err error)
}

type ServiceSync struct {
	tr ITransportSync
	sg IServiceGroups
	sd IServiceData
}

func New(tr ITransportSync, sg IServiceGroups, sd IServiceData) *ServiceSync {
	return &ServiceSync{
		tr: tr,
		sg: sg,
		sd: sd,
	}
}
