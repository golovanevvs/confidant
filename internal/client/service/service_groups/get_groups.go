package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) GetGroups(ctx context.Context, email string) (groups []model.Group, err error) {
	action := "get groups"

	groups, err = sv.rp.GetGroups(ctx, email)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return groups, nil
}
