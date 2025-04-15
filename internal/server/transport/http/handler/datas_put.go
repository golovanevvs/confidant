package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (hd *handler) DatasPut(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("add datas, url: %s, method: %s", r.URL.String(), r.Method)

	// checking the Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			customerrors.ErrContentType400,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusBadRequest)
		return
	}

	var datasBase64 []model.DataBase64
	if err := json.NewDecoder(r.Body).Decode(&datasBase64); err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			customerrors.ErrDecodeJSON400,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusBadRequest)
		return
	}

	datas := make([]model.Data, 0)
	for _, dataBase64 := range datasBase64 {
		var note model.NoteEnc
		var pass model.PassEnc
		var card model.CardEnc
		var file model.FileEnc

		switch dataBase64.DataType {
		case "note":
			note.Desc, _ = base64.StdEncoding.DecodeString(dataBase64.Note.Desc)
			note.Note, _ = base64.StdEncoding.DecodeString(dataBase64.Note.Note)
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
		data := model.Data{
			ID:         dataBase64.ID,
			IDOnClient: dataBase64.IDOnClient,
			GroupID:    dataBase64.GroupID,
			DataType:   dataBase64.DataType,
			Title:      dataBase64.Title,
			CreatedAt:  dataBase64.CreatedAt,
			Note:       note,
			Pass:       pass,
			Card:       card,
			File:       file,
		}
		hd.lg.Infof("group ID: %d, data type: %s, title: %s, note desc: %v, note note: %v", data.GroupID, data.DataType, data.Title, data.Note.Desc, data.Note.Note)
		datas = append(datas, data)
	}

	dataIDs, err := hd.sv.AddDatas(r.Context(), datas)
	if err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusInternalServerError)
		return
	}

	hd.lg.Infof("dataIDs: %v", dataIDs)

	responseJSON, err := json.Marshal(dataIDs)
	if err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			customerrors.ErrEncodeJSON500,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusInternalServerError)
		return
	}
	hd.lg.Infof("dataIDsJSON: %v", responseJSON)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseJSON))
}
