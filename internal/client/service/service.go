package service

import (
	"github.com/golovanevvs/confidant/internal/client/model"
)

type ITransport interface {
	RegisterAccount(email, password string) (trResponse *model.TrResponse, err error)
	ServerStatus() (statusResp *model.TrResponse, err error)
}

type Service struct {
	tr ITransport
}

func New(tr ITransport) *Service {
	return &Service{
		tr: tr,
	}
}
