package service_groups

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceGroups) AddEmailsBySync(ctx context.Context, mapGroupIDEmails map[int][]string) (err error) {
	action := "add e-mails by sync"
	err = sv.rp.AddEmailsBySync(ctx, mapGroupIDEmails)
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
