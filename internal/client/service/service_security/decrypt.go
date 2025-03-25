package service_security

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func (sv *ServiceSecurity) Decrypt(data []byte) (decryptedData []byte, err error) {
	action := "decrypt"

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}

	nonceSize := gcm.NonceSize()

	nonce, data := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, data, nil)
}
