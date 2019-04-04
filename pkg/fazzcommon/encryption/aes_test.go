package encryption

import (
	"testing"
)

func TestEncryptAES(t *testing.T) {
	_, err := EncryptAES([]byte("1234567890123456"), "Hello World")
	if err != nil {
		t.Fatalf("failed test encrypt aes")
	}
}

func TestFailedEncryptAES(t *testing.T) {
	_, err := EncryptAES([]byte("1"), "Hello World")
	if err == nil {
		t.Fatalf("failed test failed encrypt aes")
	}
}

func TestDescryptAES(t *testing.T) {
	_, err := DecryptAES([]byte("1234567890123456"), "-3mOzQsAHVEVT9rtNB1uXG_wObk_PzmGP6StawONRbI")
	if err != nil {
		t.Fatalf("failed test decrypt aes")
	}
}

func TestFailedDescryptAES(t *testing.T) {
	_, err := DecryptAES([]byte(""), "-3mOzQsAHVEVT9rtNB1uXG_wObk_PzmGP6StawONRbI")
	if err == nil {
		t.Fatalf("failed test failed decrypt aes")
	}
	_, err = DecryptAES([]byte("1234567890123456"), "hei )(*&^%^&*(*&^")
	if err == nil {
		t.Fatalf("failed test second failed decrypt aes")
	}
	_, err = DecryptAES([]byte("1234567890123456"), "hei")
	if err == nil {
		t.Fatalf("failed test second failed decrypt aes")
	}
}
