package service_security

import (
	"fmt"
	"os"
)

func (sv *ServiceSecurity) DecryptFile(data []byte, filepath string) (err error) {
	action := "decrypt file"

	decrypted, err := sv.Decrypt(data)
	if err != nil {
		return fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}

	return os.WriteFile(filepath, decrypted, 0644)
}
