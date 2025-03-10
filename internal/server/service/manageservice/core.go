package manageservice

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ManageService) PingDB() (err error) {
	action := "ping DB"
	if err = sv.Rp.Ping(); err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ManageServiceErr,
			action,
			err,
		)
	}
	return nil
}
