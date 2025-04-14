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
	trResponse, err := sv.tr.GetDataIDs(ctx, accessToken)
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
	trResponse.DataIDs = response.IDs

	dataIDsFromServer := trResponse.DataIDs

	// gettingdata server IDs and local IDs from client
	dataServerIDs, dataNoServerIDs, err := sv.sd.GetDataIDs(ctx, groupIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	//TODO: Add date comparison after the edit operation is added.
	// getting:
	// - data IDs for copy to client from server
	// - data IDs for comparison dates
	if len(dataIDsFromServer) > 0 {
		dataIDsForCopyFromServer := make([]int, 0)
		// dataIDsForComparisonDates := make([]int, 0)
		for _, dataIDFromServer := range dataIDsFromServer {
			if !slices.Contains(dataServerIDs, dataIDFromServer) {
				dataIDsForCopyFromServer = append(dataIDsForCopyFromServer, dataIDFromServer)
			} else {
				// dataIDsForComparisonDates = append(dataIDsForComparisonDates, dataIDFromServer)
			}
		}

		// if len(dataIDsForComparisonDates) > 0 {
		// 	dataDatesFromServer, err := sv.tr.GetDataDates(ctx, accessToken, dataIDsForComparisonDates)
		// 	if err != nil {
		// 		return nil, fmt.Errorf(
		// 			"%s: %s: %s: %w",
		// 			customerrors.ClientMsg,
		// 			customerrors.ClientServiceErr,
		// 			action,
		// 			err,
		// 		)
		// 	}

		// 	dataDatesFromClient, err := sv.sd.GetDataDates(ctx, dataIDsForComparisonDates)
		// 	if err != nil {
		// 		return nil, fmt.Errorf(
		// 			"%s: %s: %s: %w",
		// 			customerrors.ClientMsg,
		// 			customerrors.ClientServiceErr,
		// 			action,
		// 			err,
		// 		)
		// 	}

		// 	// comparison dates
		// 	// getting:
		// 	// - data IDs for update from server to client
		// 	// - data IDs for update from client to server
		// 	dataIDsForUpdateFromServerToClient := make([]int, 0)
		// 	dataIDsForUpdateFromClientToServer := make([]int, 0)
		// 	for _, dataID := range dataIDsForComparisonDates {
		// 		if dataDatesFromServer[dataID].After(dataDatesFromClient[dataID]) {
		// 			dataIDsForUpdateFromServerToClient = append(dataIDsForUpdateFromServerToClient, dataID)
		// 		}
		// 		if dataDatesFromClient[dataID].After(dataDatesFromServer[dataID]) {
		// 			dataIDsForUpdateFromClientToServer = append(dataIDsForUpdateFromClientToServer, dataID)
		// 		}
		// 	}
		// }

		if len(dataIDsForCopyFromServer) > 0 {
			// getting data from server
			datasFromServer, err := sv.tr.GetDatas(ctx, accessToken, dataIDsForCopyFromServer)
			if err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %s: %w",
					customerrors.ClientMsg,
					customerrors.ClientServiceErr,
					action,
					err,
				)
			}

			// saving datas to client DB
			err = sv.sd.SaveDatas(ctx, datasFromServer)
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

			// getting files from server
			for _, dataFromServer := range datasFromServer {
				if dataFromServer.DataType == "file" {
					file, err := sv.tr.GetDataFile(ctx, accessToken, dataFromServer.IDOnServer)
					if err != nil {
						return nil, fmt.Errorf(
							"%s: %s: %s: %w",
							customerrors.ClientMsg,
							customerrors.ClientServiceErr,
							action,
							err,
						)
					}

					// saving file to client DB
					err = sv.sd.SaveDataFile(ctx, dataFromServer.IDOnServer, file)
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

		// 		// adding group to client DB
	}

	// getting datas from client
	if len(dataNoServerIDs) > 0 {
		datasForCopyToServer, err := sv.sd.GetDatas(ctx, dataNoServerIDs)
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
	//TODO: отправить данные на сервер
	//TODO: получить от сервера id_on_server
	//TODO: обновить id_on_server у клиента

	//TODO: написать хендлер, сервис, БД на сервере с преобразованием данных в []byte

	//TODO: сделать то же самое для сохранения на сервер файла

	//! ----------------- СТОП --------------------
	// 	// sending groups to server
	// 	newGroupIDsFromServer, err := sv.tr.SendGroups(ctx, accessToken, groupsForCopyToServer)
	// 	if err != nil {
	// 		return nil, fmt.Errorf(
	// 			"%s: %s: %s: %w",
	// 			customerrors.ClientMsg,
	// 			customerrors.ClientServiceErr,
	// 			action,
	// 			err,
	// 		)
	// 	}

	// 	// updating IDsOnServer
	// 	err = sv.sg.UpdateGroupIDsOnServer(ctx, newGroupIDsFromServer)
	// 	if err != nil {
	// 		return nil, fmt.Errorf(
	// 			"%s: %s: %s: %w",
	// 			customerrors.ClientMsg,
	// 			customerrors.ClientServiceErr,
	// 			action,
	// 			err,
	// 		)
	// 	}
	// }

	// return &model.SyncResp{
	// 	HTTPStatusCode: trResponse.HTTPStatusCode,
	// 	HTTPStatus:     trResponse.HTTPStatus,
	// 	Error:          "",
	// }, nil
	return
}
