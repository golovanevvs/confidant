package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (hd *handler) EmailsPatch(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("add emails, url: %s, method: %s", r.URL.String(), r.Method)

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

	var mapGroupIDEmails map[int][]string
	if err := json.NewDecoder(r.Body).Decode(&mapGroupIDEmails); err != nil {
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

	err := hd.sv.AddEmails(r.Context(), mapGroupIDEmails)
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

	w.WriteHeader(http.StatusOK)
}
