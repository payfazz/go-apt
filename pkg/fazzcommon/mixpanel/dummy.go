package mixpanel

import "github.com/dukex/mixpanel"

type dummy struct {
}

func (d *dummy) Track(distinctId, eventName string, e *mixpanel.Event) error {
	return nil
}

func (d *dummy) Update(distinctId string, u *mixpanel.Update) error {
	return nil
}

func (d *dummy) Alias(distinctId, newId string) error {
	return nil
}

func NewDummyMixPanelClient() mixpanel.Mixpanel {
	return &dummy{}
}
