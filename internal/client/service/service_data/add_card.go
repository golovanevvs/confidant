package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) AddCard(ctx context.Context, data model.CardDec, accountID int, groupID int) (err error) {
	action := "add card"

	dataEnc := model.CardEnc{
		GroupID: groupID,
		Type:    "card",
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

	dataEnc.Number, err = sv.ss.Encrypt([]byte(data.Number))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Date, err = sv.ss.Encrypt([]byte(data.Date))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Name, err = sv.ss.Encrypt([]byte(data.Name))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.CVC2, err = sv.ss.Encrypt([]byte(data.CVC2))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.PIN, err = sv.ss.Encrypt([]byte(data.PIN))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Bank, err = sv.ss.Encrypt([]byte(data.Bank))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	err = sv.rp.AddCard(ctx, dataEnc)
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
