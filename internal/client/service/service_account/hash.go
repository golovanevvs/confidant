package service_account

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"golang.org/x/crypto/bcrypt"
)

// genHash generates a hash of the word
func (sv *ServiceAccount) genHash(word string) ([]byte, error) {
	shaHash := sha256.Sum256([]byte(word))
	shaHashHex := hex.EncodeToString(shaHash[:])

	hash, err := bcrypt.GenerateFromPassword([]byte(shaHashHex), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %w",
			customerrors.ErrGenHash,
			err,
		)
	}

	return hash, nil
}
