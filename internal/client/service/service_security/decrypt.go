package service_security

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
)

func (sv *ServiceSecurity) Decrypt(data []byte) (decryptedData []byte, err error) {
	action := "decrypt"

	if len(data) == 0 {
		return nil, fmt.Errorf(
			"%s: %w",
			action,
			customerrors.ErrDecryptEmptyBody,
		)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w",
			action,
			err,
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w",
			action,
			err,
		)
	}

	nonceSize := gcm.NonceSize()

	nonce, data := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, data, nil)
}
