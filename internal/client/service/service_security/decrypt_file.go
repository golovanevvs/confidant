package service_security

import (
	"fmt"
	"os"
)

func (sv *ServiceSecurity) DecryptFile(data []byte, filepath string) (err error) {
	action := "decrypt file"

	dataDec, err := sv.Decrypt(data)
	if err != nil {
		return fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}
	defer file.Close()

	_, err = file.Write(dataDec)
	if err != nil {
		return fmt.Errorf(
			"%s:%w",
			action,
			err,
		)
	}

	return nil
}
