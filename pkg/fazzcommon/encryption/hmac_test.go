package encryption

import "testing"

func TestCheckMACByString(t *testing.T) {
	res := CheckMACByString("88cd2108b5347d973cf39cdf9053d7dd42704876d8c9a9bd8e2d168259d3ddf7", []byte("test"), []byte("test"))
	if !res {
		t.Fatalf("failed to test hmac")
	}
}

func TestFailedCheckMACByString(t *testing.T) {
	res := CheckMACByString("asdf", []byte("test"), []byte("test"))
	if res {
		t.Fatalf("failed to test failed hmac")
	}
}
