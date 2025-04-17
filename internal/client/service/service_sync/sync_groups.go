package service_sync

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

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

	var groupIDsFromServer []int
	err = json.Unmarshal(trResponse.ResponseBody, &groupIDsFromServer)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}

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
		groupIDsForSyncEmails := make([]int, 0)
		for _, groupIDFromServer := range groupIDsFromServer {
			if !slices.Contains(groupServerIDs, groupIDFromServer) {
				groupIDsForCopyFromServer = append(groupIDsForCopyFromServer, groupIDFromServer)
			} else {
				groupIDsForSyncEmails = append(groupIDsForSyncEmails, groupIDFromServer)
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

			//TODO: add a title match check

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

		//sync e-mails
		if len(groupIDsForSyncEmails) > 0 {
			mapGroupIDEmailsFromServer, err := sv.tr.GetEmails(ctx, accessToken, groupIDsForSyncEmails)
			if err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %s: %w",
					customerrors.ClientMsg,
					customerrors.ClientServiceErr,
					action,
					err,
				)
			}

			mapGroupIDEmailsFromClient, err := sv.sg.GetEmails(ctx, groupIDsForSyncEmails)
			if err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %s: %w",
					customerrors.ClientMsg,
					customerrors.ClientServiceErr,
					action,
					err,
				)
			}

			mapGroupIDEmailsForAddToClient := make(map[int][]string)
			for groupIDFromServer, emailsFromServer := range mapGroupIDEmailsFromServer {
				emailsForAddToClient := make([]string, 0)
				for _, emailFromServer := range emailsFromServer {
					if !slices.Contains(mapGroupIDEmailsFromClient[groupIDFromServer], emailFromServer) {
						emailsForAddToClient = append(emailsForAddToClient, emailFromServer)
					}
					mapGroupIDEmailsForAddToClient[groupIDFromServer] = emailsForAddToClient
				}
			}

			mapGroupIDEmailsForAddToServer := make(map[int][]string)
			for groupIDFromCLient, emailsFromCLient := range mapGroupIDEmailsFromClient {
				emailsForAddToServer := make([]string, 0)
				for _, emailFromClient := range emailsFromCLient {
					if !slices.Contains(mapGroupIDEmailsFromServer[groupIDFromCLient], emailFromClient) {
						emailsForAddToServer = append(emailsForAddToServer, emailFromClient)
					}
					mapGroupIDEmailsForAddToServer[groupIDFromCLient] = emailsForAddToServer
				}
			}

			if len(mapGroupIDEmailsForAddToClient) > 0 {
				err = sv.sg.AddEmailsBySync(ctx, mapGroupIDEmailsForAddToClient)
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

			if len(mapGroupIDEmailsForAddToServer) > 0 {
				err = sv.tr.SendEmails(ctx, accessToken, mapGroupIDEmailsForAddToServer)
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
		//TODO: add send e-mails
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
