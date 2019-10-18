package vault

import (
	"testing"

	"github.com/payfazz/go-apt/pkg/fazzkv/env"
	"github.com/stretchr/testify/require"
)

var envClient = env.NewFazzEnv()
var client Interface

func TestFailedAuthNewFazzVault(t *testing.T) {
	var err error
	url, _ := envClient.Get("V_URL")
	client, err = NewFazzVault(url, "asdf", "asdf")
	require.Error(t, err)
}

func TestFailedNewFazzVault(t *testing.T) {
	var err error
	client, err = NewFazzVault("https://localhost:12334", "asdf", "asdf")
	require.Error(t, err)
}

func TestNewFazzVault(t *testing.T) {
	var err error
	url, _ := envClient.Get("V_URL")
	username, _ := envClient.Get("V_USERNAME")
	password, _ := envClient.Get("V_PASSWORD")
	client, err = NewFazzVault(url, username, password)
	require.NoError(t, err)
}

func TestVault_FailedGet(t *testing.T) {
	_, err := client.Get("key")
	require.Error(t, err)
}

func TestVault_ReadPath(t *testing.T) {
	var err error
	client, err = client.ReadPath("secret/fazzcard/backend")
	require.NoError(t, err)
}

func TestVault_Get(t *testing.T) {
	result, err := client.Get("key")
	require.NoError(t, err)
	require.Equal(t, "value", result)
}

func TestVault_Set(t *testing.T) {
	err := client.Set("", "")
	require.NoError(t, err)
}

func TestVault_Truncate(t *testing.T) {
	err := client.Truncate()
	require.NoError(t, err)
}

func TestVault_Delete(t *testing.T) {
	err := client.Delete("")
	require.NoError(t, err)
}
