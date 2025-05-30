package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
)

func (hd *handler) accountRegisterPost(w http.ResponseWriter, r *http.Request) {
	action := fmt.Sprintf("account register, url: %s, method: %s", r.URL.String(), r.Method)

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

	// e-mail validation
	if err := account.ValidateEmail(); err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	// password validation
	if err := account.ValidatePassword(); err != nil {
		resErr := fmt.Errorf(
			"%s: %s: %s: %w",
			customerrors.ServerMsg,
			customerrors.HandlerErr,
			action,
			err,
		)
		hd.lg.Errorf(resErr.Error())
		http.Error(w, resErr.Error(), http.StatusUnprocessableEntity)
		return
	}

	// launching the createAccount service,
	// obtaining the account ID of a new account for subsequent authorization,
	// error checking
	accountID, err := hd.sv.CreateAccount(r.Context(), account)
	if err != nil {
		switch {
		// if the email already exists in the DB
		case errors.Is(err, customerrors.ErrDBBusyEmail409):
			resErr := fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ServerMsg,
				customerrors.HandlerErr,
				action,
				err,
			)
			hd.lg.Errorf(resErr.Error())
			http.Error(w, resErr.Error(), http.StatusConflict)
			return
		// other errors
		case errors.Is(err, customerrors.ErrDBInternalError500):
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
	}

	// saving the accountId in account if there were no errors
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
	response := model.Account{
		ID: account.ID,
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
