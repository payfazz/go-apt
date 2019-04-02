package encryption

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSHA is the function that used to hashing the data with SHA method.
func HashSHA(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
