package service

import (
	"context"

	"github.com/golovanevvs/confidant/internal/server/repository"
)

type myService struct {
	myRp repository.IMyRepository
}

func newMyService(myRp repository.IMyRepository) *myService {
	return &myService{
		myRp: myRp,
	}
}

func (sv *myService) MyTask(ctx context.Context) (int, error) {
	return 0, nil
}
