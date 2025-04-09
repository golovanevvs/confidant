package service_sync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceSync) SyncGroups(ctx context.Context, accessToken string, email string) (syncResp *model.SyncResp, err error) {
	action := "sync groups"

	// getting group IDs from server
	trResponse, err := sv.tr.GetGroupIDs(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	if trResponse.HTTPStatusCode != 200 {
		return &model.SyncResp{
			HTTPStatusCode: trResponse.HTTPStatusCode,
			HTTPStatus:     trResponse.HTTPStatus,
			Error: fmt.Sprintf(
				"%s: %s: %s: %s",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				string(trResponse.ResponseBody),
			),
		}, nil
	}

	var groupIDs map[int]struct{}
	err = json.Unmarshal(trResponse.ResponseBody, &groupIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}

	groupIDsFromServer := trResponse.GroupIDs

	// getting group server IDs and local IDs from client
	groupServerIDs, groupNoServerIDs, err := sv.sg.GetGroupIDs(ctx, email)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// getting group IDs for copy to client from server
	groupIDsForCopyFromServer := make([]int, 0)
	for groupIDFromServer := range groupIDsFromServer {
		if _, inMap := groupServerIDs[groupIDFromServer]; !inMap {
			groupIDsForCopyFromServer = append(groupIDsForCopyFromServer, groupIDFromServer)
		}
	}

	// getting groups from server
	groupsFromServer, err := sv.tr.GetGroups(ctx, accessToken, groupIDsForCopyFromServer)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	//TODO: добавить проверку совпадения title

	// adding group to client DB
	for _, groupFromServer := range groupsFromServer {
		err = sv.sg.AddGroupBySync(ctx, groupFromServer)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				err,
			)
		}
	}

	// getting groups from client
	groupIDsForCopyToServer := make([]int, 0)
	for groupNoServerID := range groupNoServerIDs {
		groupIDsForCopyToServer = append(groupIDsForCopyToServer, groupNoServerID)
	}

	groupsForCopyToServer, err := sv.sg.GetGroupsByIDs(ctx, groupIDsForCopyToServer)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// sending groups to server
	newGroupIDsFromServer, err := sv.tr.SendGroups(ctx, accessToken, groupsForCopyToServer)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// updating IDsOnServer
	err = sv.sg.UpdateGroupIDsOnServer(ctx, newGroupIDsFromServer)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	return &model.SyncResp{
		HTTPStatusCode: trResponse.HTTPStatusCode,
		HTTPStatus:     trResponse.HTTPStatus,
		Error:          "",
	}, nil
}
