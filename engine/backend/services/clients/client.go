package clients

import (
	"github.com/ogiusek/relay/v2"
)

// server relay listener
// client relay listener

type ClientEmitter interface {
	// server -> client
	Relay() relay.Relay
}

type Client interface {
	ClientEmitter
	Close()
}

type ClientListener interface {
	// server <- client
	Relay() relay.Builder
	Build() Client
}

type ServerEmitter interface {
	// client -> server
	Relay() relay.Relay
}

type Server interface {
	ServerEmitter
	Close()
}

type ServerListener interface {
	// client <- server
	Relay() relay.Builder
	Build() Server
}
