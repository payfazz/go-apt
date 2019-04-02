package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// CheckMACByString is a function that used to compare string, with byte param and key param that encrypted with HMAC method.
// if string and encryption between byte and key is not the same, then it will return false. If it's same then return true.
func CheckMACByString(messageMac string, message []byte, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	if hex.EncodeToString(mac.Sum(nil)) != messageMac {
		return false
	}
	return true
}
