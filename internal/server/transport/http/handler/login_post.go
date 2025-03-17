package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (hd *handler) loginPost(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("login, url: %s, method: %s", r.URL.String(), r.Method)

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

	// deserializing JSON in account
	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
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

	// launching the login service
	accountID, err := hd.sv.Login(r.Context(), account)
	if err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		switch {
		case errors.Is(err, customerrors.ErrDBEmailNotFound401):
			http.Error(w, resErr.Error(), http.StatusUnauthorized)
			return
		case errors.Is(err, customerrors.ErrDBWrongPassword401):
			http.Error(w, resErr.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, resErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	account.ID = accountID

	// authorization
	// getting a access token string
	accessTokenString, err := hd.sv.BuildAccessJWTString(r.Context(), account.ID)
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

	// getting a refresh token string
	refreshTokenString, err := hd.sv.BuildRefreshJWTString(r.Context(), account.ID)
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

	// creating a response
	response := model.AccountRegisterResp{
		AccountID: strconv.Itoa(account.ID),
	}

	responseJSON, err := json.MarshalIndent(response, "", " ")
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
	w.Header().Set("Authorization", fmt.Sprint("Bearer ", accessTokenString))
	w.Header().Set("Refresh-Token", refreshTokenString)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseJSON))
}
