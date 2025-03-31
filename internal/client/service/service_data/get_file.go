package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetFile(ctx context.Context, dataID int) (data *model.FileDec, err error) {
	action := "get file"

	dataEnc, err := sv.rp.GetFile(ctx, dataID)
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

	filenameDec, err := sv.ss.Decrypt(dataEnc.Filename)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	filesizeDec, err := sv.ss.Decrypt(dataEnc.Filesize)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	filedateDec, err := sv.ss.Decrypt(dataEnc.Filedate)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.FileDec{
		ID:       dataEnc.ID,
		Title:    dataEnc.Title,
		Desc:     string(descDec),
		Filename: string(filenameDec),
		Filesize: string(filesizeDec),
		Filedate: string(filedateDec),
	}, nil
}
