package encryption

import (
	"crypto/sha1"
	"encoding/hex"
)

// HashSHA1 is a function that used to hash the data with SHA1 method
func HashSHA1(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
