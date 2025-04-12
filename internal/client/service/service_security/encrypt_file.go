package service_security

import (
	"fmt"
	"os"
)

func (sv *ServiceSecurity) EncryptFile(filepath string) (encryptedFile []byte, err error) {
	action := "encrypt file"

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w",
			action,
			err,
		)
	}

	encryptedFile, err = sv.Encrypt(data)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w",
			action,
			err,
		)
	}

	return encryptedFile, nil
}
