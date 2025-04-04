package service_groups

import "context"

func (sv *ServiceGroups) Sync(ctx context.Context, accessToken string) (err error) {
	action := "sync groups"

	serverGroupIDs, err := sv.tr.GetGroupIDs(ctx, accessToken)

	return nil
}
