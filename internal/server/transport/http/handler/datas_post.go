package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (hd *handler) DatasPost(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("get datas, url: %s, method: %s", r.URL.String(), r.Method)

	var err error
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

	var dataIDs []int
	if err = json.NewDecoder(r.Body).Decode(&dataIDs); err != nil {
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
	hd.lg.Infof("dataIDs: %v", dataIDs)

	if len(dataIDs) == 0 {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			customerrors.ErrDebug,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusBadRequest)
		return
	}

	var datasBase64 []model.DataBase64
	datasBase64, err = hd.sv.GetDatas(r.Context(), dataIDs)
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
	for _, dataBase64 := range datasBase64 {
		hd.lg.Infof("ID: %d, IDOnClient: %d, GroupID: %d, DataType: %s, Title: %s, note desc: %v", dataBase64.ID, dataBase64.IDOnClient, dataBase64.GroupID, dataBase64.DataType, dataBase64.Title, dataBase64.Note.Desc)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(datasBase64); err != nil {
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
}
