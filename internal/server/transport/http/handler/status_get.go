package handler

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (hd *handler) StatusGet(w http.ResponseWriter, r *http.Request) {
	action := "get server status"
	if err := hd.sv.PingDB(); err != nil {
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
