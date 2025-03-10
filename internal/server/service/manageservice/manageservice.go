package manageservice

type IManageRepository interface {
	Ping() (err error)
}

type ManageService struct {
	Rp IManageRepository
}

func New(rp IManageRepository) *ManageService {
	return &ManageService{
		Rp: rp,
	}
}
