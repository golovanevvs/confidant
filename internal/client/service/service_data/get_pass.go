package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetPass(ctx context.Context, dataID int) (data *model.PassDec, err error) {
	action := "get password"

	dataEnc, err := sv.rp.GetPass(ctx, dataID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	descDec, err := sv.ss.Decrypt(dataEnc.Desc)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	loginDec, err := sv.ss.Decrypt(dataEnc.Login)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	passDec, err := sv.ss.Decrypt(dataEnc.Pass)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.PassDec{
		ID:    dataEnc.ID,
		Title: dataEnc.Title,
		Desc:  string(descDec),
		Login: string(loginDec),
		Pass:  string(passDec),
	}, nil
}
