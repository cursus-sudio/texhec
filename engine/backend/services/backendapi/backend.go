package backendapi

import (
	"github.com/ogiusek/relay/v2"
)

type Backend interface {
	Relay() relay.Relay
}

type backend struct {
	relay relay.Relay
}

func newBackend(relay relay.Relay) Backend {
	return &backend{relay: relay}
}

func (backend *backend) Relay() relay.Relay {
	return backend.relay
}
