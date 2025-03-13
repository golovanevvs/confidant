package service_account

import (
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"golang.org/x/crypto/bcrypt"
)

// genHash generates a hash of the word
func (sv *ServiceAccount) genHash(word string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(word), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %w",
			customerrors.ErrGenHash,
			err,
		)
	}

	return hash, nil
}
