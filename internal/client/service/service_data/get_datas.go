package service_data

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceData) GetDatas(ctx context.Context, dataIDs []int) (datasBase64 []model.DataBase64, err error) {
	action := "get datas"

	var datas []model.Data
	datas, err = sv.rp.GetDatasByIDs(ctx, dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ClientMsg,
			customerrors.ClientServiceErr,
			action,
			err,
		)
	}

	datasBase64 = make([]model.DataBase64, len(datas))

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

	return datasBase64, nil
}
