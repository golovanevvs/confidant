package service_manage

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceManage) PingDB() (err error) {
	action := "ping DB"
	if err = sv.rp.Ping(); err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ManageServiceErr,
			action,
			err,
		)
	}
	return nil
}
