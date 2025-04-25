package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (hd *handler) DataFilePut(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("add file, url: %s, method: %s", r.URL.String(), r.Method)

	dataIDStr := r.Header.Get("X-Data-ID")
	if dataIDStr == "" {
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

	dataID, err := strconv.Atoi(dataIDStr)
	if err != nil {
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

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/octet-stream" {
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

	file, err := io.ReadAll(r.Body)
	if err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusBadRequest)
		return
	}

	err = hd.sv.AddDataFile(r.Context(), dataID, file)
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
