package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) AddGroupBySync(ctx context.Context, group model.Group) (err error) {
	action := "add group by sync"
	err = sv.rp.AddGroupBySync(ctx, group)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}
	return nil
}
