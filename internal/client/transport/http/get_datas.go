package trhttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (tr *trHTTP) GetDatas(ctx context.Context, accessToken string, dataIDs []int) (datasFromServer []model.Data, err error) {
	action := "get datas"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/datas", tr.addr)

	dataIDsJSON, err := json.Marshal(dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(dataIDsJSON))
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrCreateRequest,
			err,
		)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	//! Response
	response, err := tr.cl.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrSendRequest,
			err,
		)
	}
	defer response.Body.Close()

	var datasBase64 []model.DataBase64
	if err := json.NewDecoder(response.Body).Decode(&datasBase64); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}
	datasFromServer = make([]model.Data, len(datasBase64))
	for _, dataBase64 := range datasBase64 {
		var note model.NoteEnc
		var pass model.PassEnc
		var card model.CardEnc
		var file model.FileEnc

		switch dataBase64.DataType {
		case "note":
			note.Desc, _ = base64.StdEncoding.DecodeString(dataBase64.Note.Desc)
			note.Note, _ = base64.StdEncoding.DecodeString(dataBase64.Note.Desc)
		case "pass":
			pass.Desc, _ = base64.StdEncoding.DecodeString(dataBase64.Pass.Desc)
			pass.Login, _ = base64.StdEncoding.DecodeString(dataBase64.Pass.Login)
			pass.Pass, _ = base64.StdEncoding.DecodeString(dataBase64.Pass.Pass)
		case "card":
			card.Desc, _ = base64.StdEncoding.DecodeString(dataBase64.Card.Desc)
			card.Number, _ = base64.StdEncoding.DecodeString(dataBase64.Card.Number)
			card.Date, _ = base64.StdEncoding.DecodeString(dataBase64.Card.Date)
			card.Name, _ = base64.StdEncoding.DecodeString(dataBase64.Card.Name)
			card.CVC2, _ = base64.StdEncoding.DecodeString(dataBase64.Card.CVC2)
			card.PIN, _ = base64.StdEncoding.DecodeString(dataBase64.Card.PIN)
			card.Bank, _ = base64.StdEncoding.DecodeString(dataBase64.Card.Bank)
		case "file":
			file.Desc, _ = base64.StdEncoding.DecodeString(dataBase64.File.Desc)
			file.Filename, _ = base64.StdEncoding.DecodeString(dataBase64.File.Filename)
			file.Filesize, _ = base64.StdEncoding.DecodeString(dataBase64.File.Filesize)
			file.Filedate, _ = base64.StdEncoding.DecodeString(dataBase64.File.Filedate)
		}
		dataFromServer := model.Data{
			ID:        dataBase64.ID,
			GroupID:   dataBase64.GroupID,
			DataType:  dataBase64.DataType,
			Title:     dataBase64.Title,
			CreatedAt: dataBase64.CreatedAt,
			Note:      note,
			Pass:      pass,
			Card:      card,
			File:      file,
		}
		datasFromServer = append(datasFromServer, dataFromServer)
	}

	return datasFromServer, nil
}
