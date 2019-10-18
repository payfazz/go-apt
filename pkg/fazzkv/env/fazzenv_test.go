package env

import (
	"testing"

	"github.com/payfazz/go-apt/pkg/fazzkv"
	"github.com/stretchr/testify/require"
)

var osClient fazzkv.Store

func TestNewFazzEnv(t *testing.T) {
	osClient = NewFazzEnv()
}

func TestEnv_Set(t *testing.T) {
	err := osClient.Set("test", "test")
	require.NoError(t, err)
}

func TestEnv_Get(t *testing.T) {
	result, err := osClient.Get("test")
	require.NoError(t, err)
	require.Equal(t, "test", result)
}

func TestEnv_Delete(t *testing.T) {
	_ = osClient.Set("test2", "test2")
	err := osClient.Delete("test")
	require.NoError(t, err)
}

func TestEnv_Truncate(t *testing.T) {
	err := osClient.Truncate()
	require.NoError(t, err)
}
