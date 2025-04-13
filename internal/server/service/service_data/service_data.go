package service_data

import (
	"context"
	"time"
)

type IRepositoryData interface {
	GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error)
	GetDataDates(ctx context.Context, dataIDs []int) (datadates map[int]time.Time, err error)
}

type ServiceData struct {
	rp IRepositoryData
}

func New(rp IRepositoryData) *ServiceData {
	return &ServiceData{
		rp: rp,
	}
}
