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

func (tr *trHTTP) SendDatas(ctx context.Context, accessToken string, datas []model.Data) (dataIDs map[int]int, err error) {
	action := "send datas"

	//! Request
	endpoint := fmt.Sprintf("http://%s/api/datas", tr.addr)

	datasBase64 := make([]model.DataBase64, 0)

	for _, data := range datas {
		var noteBase64 model.NoteBase64
		var passBase64 model.PassBase64
		var cardBase64 model.CardBase64
		var fileBase64 model.FileBase64
		var dataBase64 model.DataBase64

		switch data.DataType {
		case "note":
			noteBase64 = model.NoteBase64{
				Desc: base64.StdEncoding.EncodeToString(data.Note.Desc),
				Note: base64.StdEncoding.EncodeToString(data.Note.Note),
			}
		case "pass":
			passBase64 = model.PassBase64{
				Desc:  base64.StdEncoding.EncodeToString(data.Pass.Desc),
				Login: base64.StdEncoding.EncodeToString(data.Pass.Login),
				Pass:  base64.StdEncoding.EncodeToString(data.Pass.Pass),
			}
		case "card":
			cardBase64 = model.CardBase64{
				Desc:   base64.StdEncoding.EncodeToString(data.Card.Desc),
				Number: base64.StdEncoding.EncodeToString(data.Card.Number),
				Date:   base64.StdEncoding.EncodeToString(data.Card.Date),
				Name:   base64.StdEncoding.EncodeToString(data.Card.Name),
				CVC2:   base64.StdEncoding.EncodeToString(data.Card.CVC2),
				PIN:    base64.StdEncoding.EncodeToString(data.Card.PIN),
				Bank:   base64.StdEncoding.EncodeToString(data.Card.Bank),
			}
		case "file":
			fileBase64 = model.FileBase64{
				Desc:     base64.StdEncoding.EncodeToString(data.File.Desc),
				Filename: base64.StdEncoding.EncodeToString(data.File.Filename),
				Filesize: base64.StdEncoding.EncodeToString(data.File.Filesize),
				Filedate: base64.StdEncoding.EncodeToString(data.File.Filedate),
			}
		}

		dataBase64 = model.DataBase64{
			ID:        data.ID,
			GroupID:   data.GroupID,
			DataType:  data.DataType,
			Title:     data.Title,
			CreatedAt: data.CreatedAt,
			Note:      noteBase64,
			Pass:      passBase64,
			Card:      cardBase64,
			File:      fileBase64,
		}

		datasBase64 = append(datasBase64, dataBase64)
	}

	datasBase64JSON, err := json.Marshal(datasBase64)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
	}

	request, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(datasBase64JSON))
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

	err = json.NewDecoder(response.Body).Decode(&dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.ClientHTTPErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
	}

	return dataIDs, nil
}
