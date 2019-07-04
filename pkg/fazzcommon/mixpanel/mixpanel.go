package mixpanel

import (
	"github.com/dukex/mixpanel"
)

// Client is an interface for mix panel
type Client interface {
	GetClient() mixpanel.Mixpanel
	BulkSet(properties map[string]interface{})
	Set(key string, value interface{})
	Send(distinctID string, eventID string) error
	SendBulk(distinctID string, eventID string, properties map[string]interface{}) error
}

type mixPanelClient struct {
	client     mixpanel.Mixpanel
	properties map[string]interface{}
}

// GetClient is a function that used to get client
func (m *mixPanelClient) GetClient() mixpanel.Mixpanel {
	return m.client
}

// BulkSet is a function that used to set properties with bulk method (reset all bulk)
func (m *mixPanelClient) BulkSet(properties map[string]interface{}) {
	m.properties = properties
}

// Set is a function that used to set the data
func (m *mixPanelClient) Set(key string, value interface{}) {
	m.properties[key] = value
}

// Send is a function that used to send the data
func (m *mixPanelClient) Send(distinctID string, eventID string) error {
	return m.client.Track(distinctID, eventID, &mixpanel.Event{
		Properties: m.properties,
	})
}

// SendBulk is a function that used to send the data
func (m *mixPanelClient) SendBulk(distinctID string, eventID string, properties map[string]interface{}) error {
	return m.client.Track(distinctID, eventID, &mixpanel.Event{
		Properties: properties,
	})
}

// NewMixPanelClient is a function that used as a constructor to construct mix panel client
func NewMixPanelClient(secret string) Client {
	return &mixPanelClient{
		client: mixpanel.New(secret, "https://api.mixpanel.com"),
	}
}

// NewMixPanelClientDummy is a function that used as a constructor to construct mix panel client
func NewMixPanelClientDummy(dummy mixpanel.Mixpanel) Client {
	return &mixPanelClient{
		client:     dummy,
		properties: make(map[string]interface{}),
	}
}
