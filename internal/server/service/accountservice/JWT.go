package accountservice

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

const (
	tokenExp  = time.Hour * 3 // token lifetime
	secretKey = "sskey"       // the secret key of the token
)

// the structure of the claims
type claims struct {
	jwt.RegisteredClaims
	AccountID int
}

// BuildJWTString creates a token and returns it as a string
func (sv *accountService) BuildJWTString(ctx context.Context, login, password string) (string, error) {
	// getting the account ID from the database
	accountID, err := sv.accountRp.LoadAccountID(ctx, login, sv.genPasswordHash(password))
	if err != nil {
		// if the login/password pair is incorrect
		if strings.Contains(err.Error(), "no rows in result set") {
			return "", fmt.Errorf("%s: %s", customerrors.DBInvalidLoginPassword401, err.Error())
		}
		// if there is another error
		return "", fmt.Errorf("%s: %s", customerrors.InternalServerError500, err.Error())
	}

	// creating a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		AccountID: accountID,
	})

	// creating a token string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %s", customerrors.InternalServerError500, err.Error())
	}

	return tokenString, nil
}

// GetAccountIDFromJWT returns the accountID from JWT
func (as *accountService) GetAccountIDFromJWT(tokenString string) (int, error) {
	claims := &claims{}

	// converting a string to a token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signature method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return -1, fmt.Errorf("%s: %s", customerrors.JWTParseError401, err.Error())
	}

	// token validation
	if !token.Valid {
		return -1, errors.New(customerrors.JWTInvalidToken401)
	}

	return claims.AccountID, nil
}
