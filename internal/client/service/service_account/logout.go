package service_account

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceAccount) Logout(ctx context.Context) (err error) {
	action := "logout"

	err = sv.rp.DeleteActiveAccount(ctx)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			customerrors.ErrLogout,
			err,
		)
	}
	return nil
}
