package accountservice

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	hashKey = "asd52v04fgt2"
)

// GenPasswordHash generates a hash of the password
func (as *accountService) genPasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write([]byte(hashKey))
	return hex.EncodeToString(hash.Sum(nil))
}
