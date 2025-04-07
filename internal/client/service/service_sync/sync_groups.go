package service_sync

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceSync) SyncGroups(ctx context.Context, accessToken string, email string) (err error) {
	action := "sync groups"

	// getting group IDs from server
	groupIDsFromServer, err := sv.tr.GetGroupIDs(ctx, accessToken)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// getting group server IDs and local IDs from client
	groupServerIDs, groupNoServerIDs, err := sv.sg.GetGroupIDs(ctx, email)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// getting group IDs for copy to client from server
	groupIDsForCopyFromServer := make(map[int]struct{})
	for groupIDFromServer := range groupIDsFromServer {
		if _, inMap := groupServerIDs[groupIDFromServer]; !inMap {
			groupIDsForCopyFromServer[groupIDFromServer] = struct{}{}
		}
	}

	// getting groups from server
	groupsFromServer, err := sv.tr.GetGroups(ctx, accessToken, groupIDsForCopyFromServer)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// adding group to client DB
	for _, groupFromServer := range groupsFromServer {
		err = sv.sg.AddGroupBySync(ctx, groupFromServer)
		if err != nil {
			return fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				err,
			)
		}
	}

	fmt.Println(groupServerIDs, groupNoServerIDs)
	return nil
}
