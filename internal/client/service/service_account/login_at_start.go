package service_account

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceAccount) LoginAtStart(ctx context.Context) (accountID int, email string, refreshTokenString string, err error) {
	action := "login at start"

	accountID, refreshTokenString, err = sv.rp.LoadActiveAccount(ctx)
	if err != nil {
		return -1, "", "", fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			customerrors.ErrLoadActiveAccount,
			err,
		)
	}

	email, err = sv.rp.LoadEmail(ctx, accountID)
	if err != nil {
		return -1, "", "", fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			customerrors.ErrLoadEmail,
			err,
		)
	}

	return accountID, email, refreshTokenString, nil
}
