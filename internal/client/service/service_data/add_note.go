package service_data

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) AddNote(ctx context.Context, data *model.NoteDec, accountID int, titleGroup string) (err error) {
	action := "add note"

	groupID, err := sv.rp.GetGroupID(ctx, accountID, titleGroup)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	dataEnc := &model.NoteEnc{
		GroupID: groupID,
	}

	dataEnc.Title, err = sv.ss.Encrypt([]byte(data.Title))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
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

	dataEnc.Note, err = sv.ss.Encrypt([]byte(data.Note))
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	err = sv.rp.AddNote(ctx, dataEnc)
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
