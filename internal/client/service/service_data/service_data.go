package service_data

type ITransportData interface {
}

type IRepositoryData interface {
}

type ServiceData struct {
	tr ITransportData
	rp IRepositoryData
}

func New(tr ITransportData, rp IRepositoryData) *ServiceData {
	return &ServiceData{
		tr: tr,
		rp: rp,
	}
}
