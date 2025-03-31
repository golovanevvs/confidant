package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetCard(ctx context.Context, dataID int) (data *model.CardDec, err error) {
	action := "get card"

	dataEnc, err := sv.rp.GetCard(ctx, dataID)
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

	numberDec, err := sv.ss.Decrypt(dataEnc.Number)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dateDec, err := sv.ss.Decrypt(dataEnc.Date)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}
	nameDec, err := sv.ss.Decrypt(dataEnc.Name)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	cvc2Dec, err := sv.ss.Decrypt(dataEnc.CVC2)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	pinDec, err := sv.ss.Decrypt(dataEnc.PIN)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	bankDec, err := sv.ss.Decrypt(dataEnc.Bank)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.CardDec{
		ID:     dataEnc.ID,
		Title:  dataEnc.Title,
		Desc:   string(descDec),
		Number: string(numberDec),
		Date:   string(dateDec),
		Name:   string(nameDec),
		CVC2:   string(cvc2Dec),
		PIN:    string(pinDec),
		Bank:   string(bankDec),
	}, nil
}
