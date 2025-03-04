package service

import "net/http"

type ITransport interface {
	RegisterAccount(email, password string) (response *http.Response, err error)
}

type Service struct {
	tr ITransport
}

func New(tr ITransport) *Service {
	return &Service{
		tr: tr,
	}
}
