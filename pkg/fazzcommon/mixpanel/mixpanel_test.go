package mixpanel

import (
	"testing"

	"github.com/dukex/mixpanel"
	"github.com/stretchr/testify/require"
)

var dummyClient = NewDummyMixPanelClient()
var mpService = NewMixPanelClientDummy(dummyClient)
var trueMPService = NewMixPanelClient("test")

func TestMixPanelClient_Set(t *testing.T) {
	mpService.Set("test", "test")
}

func TestMixPanelClient_BulkSet(t *testing.T) {
	mpService.BulkSet(map[string]interface{}{
		"test":   "test",
		"test12": "test12",
	})
}

func TestMixPanelClient_Send(t *testing.T) {
	err := mpService.Send("test", "test")
	require.NoError(t, err)
}

func TestMixPanelClient_SendBulk(t *testing.T) {
	err := mpService.SendBulk("test", "test", map[string]interface{}{
		"test":   "test",
		"test12": "test12",
	})
	require.NoError(t, err)
}

func TestMixPanelClient_GetClient(t *testing.T) {
	client := mpService.GetClient()
	err := client.Alias("test", "test")
	require.NoError(t, err)

	err = client.Update("test", &mixpanel.Update{})
	require.NoError(t, err)
}
