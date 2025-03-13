package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

type contextKey string

const AccountIDContextKey contextKey = "accountID"

// user authorization using the Authorization header token
func (hd *handler) authByJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action := "middleware authorization"

		header := r.Header.Get("Authorization")

		if header == "" {
			resErr := fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ServerMsg,
				customerrors.HandlerErr,
				action,
				customerrors.ErrAuthHeaderReq,
			)
			hd.lg.Errorf(resErr.Error())
			http.Error(w, resErr.Error(), http.StatusUnauthorized)
			return
		}

		headerSplit := strings.Split(header, " ")

		if len(headerSplit) != 2 {
			resErr := fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ServerMsg,
				customerrors.HandlerErr,
				action,
				customerrors.ErrInvalidAuthHeaderReq,
			)
			hd.lg.Errorf(resErr.Error())
			http.Error(w, resErr.Error(), http.StatusUnauthorized)
			return
		}

		if headerSplit[0] != "Bearer" {
			resErr := fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ServerMsg,
				customerrors.HandlerErr,
				action,
				customerrors.ErrBearer,
			)
			hd.lg.Errorf(resErr.Error())
			http.Error(w, resErr.Error(), http.StatusUnauthorized)
			return
		}

		if len(headerSplit[1]) == 0 {
			resErr := fmt.Errorf(
				"%s: %s: %s: %w",
				customerrors.ServerMsg,
				customerrors.HandlerErr,
				action,
				customerrors.ErrAccessToken,
			)
			hd.lg.Errorf(resErr.Error())
			http.Error(w, resErr.Error(), http.StatusUnauthorized)
			return
		}

		accountID, err := hd.sv.GetAccountIDFromJWT(headerSplit[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AccountIDContextKey, accountID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
