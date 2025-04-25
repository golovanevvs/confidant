package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) AddGroupsBySync(ctx context.Context, groups []model.Group) (err error) {
	action := "add group by sync"
	err = sv.rp.AddGroupsBySync(ctx, groups)
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
