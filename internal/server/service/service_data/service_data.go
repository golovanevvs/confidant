package service_data

import (
	"context"
	"time"

	"github.com/golovanevvs/confidant/internal/server/model"
)

type IRepositoryData interface {
	GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error)
	GetDataDates(ctx context.Context, dataIDs []int) (mapDataIDDate map[int]time.Time, err error)
	GetDatas(ctx context.Context, dataIDs []int) (datas []model.Data, err error)
	GetDataFile(ctx context.Context, dataID int) (file []byte, err error)
	AddDatas(ctx context.Context, datas []model.Data) (dataIDs map[int]int, err error)
	SaveDataFile(ctx context.Context, dataID int, file []byte) (err error)
}

type ServiceData struct {
	rp IRepositoryData
}

func New(rp IRepositoryData) *ServiceData {
	return &ServiceData{
		rp: rp,
	}
}
