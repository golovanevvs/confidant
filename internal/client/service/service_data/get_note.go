package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetNote(ctx context.Context, dataID int) (data *model.NoteDec, err error) {
	action := "get note"

	dataEnc, err := sv.rp.GetNote(ctx, dataID)
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

	noteDec, err := sv.ss.Decrypt(dataEnc.Note)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.NoteDec{
		ID:    dataEnc.ID,
		Title: dataEnc.Title,
		Desc:  string(descDec),
		Note:  string(noteDec),
	}, nil
}
