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
	// checking the Content-Type
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
	default:
		http.Error(w, string(customerrors.InvalidContentType400), http.StatusBadRequest)
		return
	}

	// deserializing JSON in Account
	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, string(customerrors.DecodeJSONError500), http.StatusInternalServerError)
		return
	}

	// e-mail validation
	// if !model.Account. (account.Email) {
	// 	return -1, errors.New("e-mail validation error")
	// }

	// launching the createUser service,
	// obtaining the account ID of a new user for subsequent authorization,
	// error checking
	accountID, err := hd.sv.IAccountService.CreateAccount(r.Context(), account)
	if err != nil {
		switch {
		case errors.Is(err, customerrors.ErrDBBusyEmail409):
			// if the email already exists in the DB
			http.Error(w, err.Error(), http.StatusConflict)
			return
			// other errors
		case errors.Is(err, customerrors.ErrDBInternalError500):
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// saving the accountId in account if there were no errors
	account.ID = accountID

	// authorization
	// getting a token string
	tokenString, err := hd.sv.IAccountService.BuildJWTString(r.Context(), account.Email, account.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// forming a response
	resMap := make(map[string]interface{})
	resMap["email"] = account.Email
	resMap["accountID"] = account.ID
	resMap["token"] = tokenString

	res, err := json.MarshalIndent(resMap, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// writing the headers and response
	w.Header().Add("Authorization", fmt.Sprint("Bearer ", tokenString))
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
