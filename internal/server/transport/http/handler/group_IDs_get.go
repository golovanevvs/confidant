package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (hd *handler) GroupIDsGet(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("get group IDs, url: %s, method: %s", r.URL.String(), r.Method)

	accountID := r.Context().Value(AccountIDContextKey).(int)
	groupIDs, err := hd.sv.GetGroupIDs(r.Context(), accountID)
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

	responseJSON, err := json.Marshal(groupIDs)
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

	// writing the headers and response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
