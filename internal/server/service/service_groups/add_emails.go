package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) AddEmails(ctx context.Context, mapGroupIDEmails map[int][]string) (err error) {
	action := "add e-mails"

	err = sv.rp.AddEmails(ctx, mapGroupIDEmails)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.GroupsServiceErr,
			action,
			err,
		)
	}

	return nil
}
