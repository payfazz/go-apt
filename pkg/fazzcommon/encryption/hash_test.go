package encryption

import (
	"testing"
)

func TestHashBCrypt(t *testing.T) {
	_, err := HashBCrypt("test")
	if err != nil {
		t.Fatalf("failed to test hash bcrypt")
	}
}

func TestCheckBCrypt(t *testing.T) {
	pass, _ := HashBCrypt("test")
	res := CheckBCrypt("test", pass)
	if !res {
		t.Fatalf("failed to test check has bcrypt")
	}
}
