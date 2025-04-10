package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (hd *handler) refreshAccessTokenPost(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("refresh access token, url: %s, method: %s", r.URL.String(), r.Method)

	// checking the Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
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

	refreshToken, err := io.ReadAll(r.Body)
	if err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusUnauthorized)
		return
	}
	defer r.Body.Close()

	// launching the service
	accessTokenString, err := hd.sv.RefreshAccessJWT(r.Context(), string(refreshToken))
	if err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusUnauthorized)
		return
	}

	// writing the headers and response
	w.Header().Set("Authorization", fmt.Sprint("Bearer ", accessTokenString))
	w.WriteHeader(http.StatusOK)
}
