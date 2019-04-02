package encryption

import "golang.org/x/crypto/bcrypt"

// HashBCrypt is a function that used to encrypt some string with BCrypt method.
func HashBCrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckBCrypt is a function that used to check some string with hashed string with BCrypt method.
func CheckBCrypt(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
