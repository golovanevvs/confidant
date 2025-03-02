package otherservice

type IManageRepository interface {
	CloseDB() error
}

type OtherService struct {
	Rp IManageRepository
}

func New(manageRp IManageRepository) *OtherService {
	return &OtherService{
		Rp: manageRp,
	}
}

func (sv *OtherService) DoSomething() error {
	return nil
}
