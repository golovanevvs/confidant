package accountservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golovanevvs/confidant/internal/customerrors"
)

const (
	accessTokenExp  = time.Hour * 3  // access token lifetime
	refreshTokenExp = time.Hour * 24 // refresh token lifetime
	secretKey       = "sskey"        // the secret key of the token
)

// the structure of the claims
type claims struct {
	jwt.RegisteredClaims
	AccountID int
}

// BuildAccessJWTString creates a access token and returns it as a string
func (sv *AccountService) BuildAccessJWTString(ctx context.Context, accountID int) (accessTokenString string, err error) {
	action := "build JWT string"
	// creating a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExp)),
		},
		AccountID: accountID,
	})

	// creating a token string
	accessTokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w: %w", customerrors.AccountServiceErr, action, customerrors.ErrTokenSignedString, err)
	}

	return accessTokenString, nil
}

// BuildRefreshJWTString creates a refresh token and returns it as a string
func (sv *AccountService) BuildRefreshJWTString(ctx context.Context, accountID int) (refreshTokenString string, err error) {
	action := "build JWT string"
	// creating a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExp)),
		},
		AccountID: accountID,
	})

	// creating a token string
	refreshTokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w: %w", customerrors.AccountServiceErr, action, customerrors.ErrTokenSignedString, err)
	}

	return refreshTokenString, nil
}

// GetAccountIDFromJWT returns the accountID from JWT
// func (sv *AccountService) GetAccountIDFromJWT(tokenString string) (int, error) {
// 	action := "get account ID from JWT"

// 	claims := &claims{}

// 	// converting a string to a token
// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("%s: %s: %w", customerrors.AccountServiceErr, action, customerrors.ErrJWTWrongSingingMethod401)
// 		}
// 		return []byte(secretKey), nil
// 	})
// 	if err != nil {
// 		return -1, fmt.Errorf("%s: %s: %w: %w", customerrors.AccountServiceErr, action, customerrors.ErrJWTInvalidToken401, err)
// 	}

// 	// token validation
// 	if !token.Valid {
// 		return -1, fmt.Errorf("%s: %s: %w", customerrors.AccountServiceErr, action, customerrors.ErrJWTInvalidToken401)
// 	}

// 	return claims.AccountID, nil
// }
