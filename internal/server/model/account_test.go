package model

import (
	"testing"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/stretchr/testify/assert"
)

func TestAccountValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  error
	}{
		{
			name:  "valid email",
			email: "test@test.com",
			want:  nil,
		},
		{
			name:  "invalid email: without @",
			email: "testtest.com",
			want:  customerrors.ErrAccountValidateEmail,
		},
		{
			name:  "invalid email: without .",
			email: "test@testcom",
			want:  customerrors.ErrAccountValidateEmail,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := &Account{
				Email: test.email,
			}
			res := account.ValidateEmail()
			assert.ErrorIs(t, res, test.want)
		})
	}
}

func TestAccountValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     error
	}{
		{
			name:     "valid password",
			password: "FgF.sd23",
			want:     nil,
		},
		{
			name:     "invalid password: len < 8",
			password: "FgF.sd2",
			want:     customerrors.ErrAccountValidatePassword,
		},
		{
			name:     "invalid password: not latin",
			password: "Fgф.sd23",
			want:     customerrors.ErrAccountValidatePassword,
		},
		{
			name:     "invalid password: without upper",
			password: "fgf.sd23",
			want:     customerrors.ErrAccountValidatePassword,
		},
		{
			name:     "invalid password: without lower",
			password: "FGF.SD23",
			want:     customerrors.ErrAccountValidatePassword,
		},
		{
			name:     "invalid password: without digit",
			password: "FgF.sdpp",
			want:     customerrors.ErrAccountValidatePassword,
		},
		{
			name:     "invalid password: without symbol",
			password: "FgF1sd23",
			want:     customerrors.ErrAccountValidatePassword,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := &Account{
				Password: test.password,
			}
			res := account.ValidatePassword()
			assert.ErrorIs(t, res, test.want)
		})
	}
}
