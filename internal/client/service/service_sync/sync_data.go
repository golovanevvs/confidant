package service_sync

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceSync) SyncData(ctx context.Context, accessToken string, email string) (syncResp *model.SyncResp, err error) {
	action := "sync data"

	// getting group IDs from client
	groupIDs, _, err := sv.sg.GetGroupIDs(ctx, email)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	// getting data IDs from server

	for _, groupID := range groupIDs {

	}

	//! ----------------- СТОП --------------------
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

	response := struct {
		IDs []int `json:"ids"`
	}{}
	err = json.Unmarshal(trResponse.ResponseBody, &response)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}
	trResponse.GroupIDs = response.IDs

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
	if len(groupIDsFromServer) > 0 {
		groupIDsForCopyFromServer := make([]int, 0)
		for _, groupIDFromServer := range groupIDsFromServer {
			if !slices.Contains(groupServerIDs, groupIDFromServer) {
				groupIDsForCopyFromServer = append(groupIDsForCopyFromServer, groupIDFromServer)
			}
		}

		if len(groupIDsForCopyFromServer) > 0 {
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
		}
	}

	// getting groups from client
	if len(groupNoServerIDs) > 0 {
		groupsForCopyToServer, err := sv.sg.GetGroupsByIDs(ctx, groupNoServerIDs)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %s: %w: groupNoServerIDs: %v",
				customerrors.ClientMsg,
				customerrors.ClientServiceErr,
				action,
				err,
				groupNoServerIDs,
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
	}

	return &model.SyncResp{
		HTTPStatusCode: trResponse.HTTPStatusCode,
		HTTPStatus:     trResponse.HTTPStatus,
		Error:          "",
	}, nil
}
