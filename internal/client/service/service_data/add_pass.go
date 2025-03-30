package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) AddPass(ctx context.Context, data model.PassDec, accountID int, groupID int) (err error) {
	action := "add password"

	dataEnc := model.PassEnc{
		GroupID: groupID,
		Type:    "pass",
		Title:   data.Title,
	}

	dataEnc.Desc, err = sv.ss.Encrypt([]byte(data.Desc))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Login, err = sv.ss.Encrypt([]byte(data.Login))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Pass, err = sv.ss.Encrypt([]byte(data.Pass))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	err = sv.rp.AddPass(ctx, dataEnc)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return nil
}
