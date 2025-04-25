package service_account

import (
	"crypto/sha256"
)

const (
	hashKey = "asd52v04fgt250"
)

// genHash generates a hash of the word
func (sv *ServiceAccount) genHash(word string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(word))
	hash.Write([]byte(hashKey))

	return hash.Sum(nil), nil
}
