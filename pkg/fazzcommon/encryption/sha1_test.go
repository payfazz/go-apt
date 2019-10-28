package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashSHA1(t *testing.T) {
	result := HashSHA1("test123")
	require.Equal(t, "7288edd0fc3ffcbe93a0cf06e3568e28521687bc", result)
}
