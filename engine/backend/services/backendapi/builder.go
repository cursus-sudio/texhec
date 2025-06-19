package backendapi

import "github.com/ogiusek/relay/v2"

type Builder interface {
	Relay() relay.Builder
	Build() Backend
}

type builder struct {
	relay relay.Builder
}

func NewBuilder() Builder {
	return &builder{
		relay: relay.NewBuilder(),
	}
}

func (backendBuilder builder) Relay() relay.Builder {
	return backendBuilder.relay
}

func (backendBuilder builder) Build() Backend {
	return newBackend(
		backendBuilder.relay.Build(),
	)
}
