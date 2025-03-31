package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) AddFile(ctx context.Context, data model.FileDec, accountID int, groupID int, filepath string) (err error) {
	action := "add file"

	dataEnc := model.FileEnc{
		GroupID: groupID,
		Type:    "file",
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

	dataEnc.Filename, err = sv.ss.Encrypt([]byte(data.Filename))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Filesize, err = sv.ss.Encrypt([]byte(data.Filesize))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.Filedate, err = sv.ss.Encrypt([]byte(data.Filedate))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc.File, err = sv.ss.EncryptFile(filepath)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	err = sv.rp.AddFile(ctx, dataEnc)
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
