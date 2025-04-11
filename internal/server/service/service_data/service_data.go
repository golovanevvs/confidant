package service_data

import (
	"context"
)

type IRepositoryData interface {
	GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error)
}

type ServiceData struct {
	rp IRepositoryData
}

func New(rp IRepositoryData) *ServiceData {
	return &ServiceData{
		rp: rp,
	}
}
