package service_manage

type IRepositoryManage interface {
	Ping() (err error)
}

type ServiceManage struct {
	rp IRepositoryManage
}

func New(rp IRepositoryManage) *ServiceManage {
	return &ServiceManage{
		rp: rp,
	}
}
