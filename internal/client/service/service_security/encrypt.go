package service_security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

var key = []byte("lakjsghpoiurKJwgo2435-23n")

func (sv *ServiceSecurity) Encrypt(data []byte) (encryptedData []byte, err error) {
	action := "encrypt"
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}

	encryptedData = gcm.Seal(nonce, nonce, data, nil)

	return encryptedData, nil
}
