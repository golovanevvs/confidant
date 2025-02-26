package model

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

type Account struct {
	ID           int    `json:"-"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
}

func (m *Account) ValidateEmail() error {
	action := "validate e-mail"

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(emailPattern)
	if !regex.MatchString(m.Email) {
		return fmt.Errorf("%s: %s: %w", customerrors.AccountErr, action, customerrors.ErrAccountValidateEmail422)
	}

	return nil
}

func (m *Account) ValidatePassword() error {
	action := "validate password"

	passwordPattern := `^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]*$`
	regex := regexp.MustCompile(passwordPattern)
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	switch {
	case len(m.Password) < 8:
		return fmt.Errorf("%s: %s: %w", customerrors.AccountErr, action, customerrors.ErrAccountValidatePassword422)
	case !regex.MatchString(m.Password):
		return fmt.Errorf("%s: %s: %w", customerrors.AccountErr, action, customerrors.ErrAccountValidatePassword422)
	default:
		for _, char := range m.Password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsDigit(char):
				hasDigit = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}
		if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
			return fmt.Errorf("%s: %s: %w", customerrors.AccountErr, action, customerrors.ErrAccountValidatePassword422)
		}
	}

	return nil
}
