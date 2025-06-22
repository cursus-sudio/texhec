package api

import (
	"github.com/ogiusek/relay/v2"
)

// server uses client (CLientEmitter, Client and ClientListener is used by server)
// client uses server (ServerEmitter, Server and ServerListener is used by client)

// services used by server

type Client interface { // server -> client
	Relay() relay.Relay
	Close()
}

type ClientBuilder interface { // server <- client
	Relay() relay.Builder
	Build(close func()) Client
}

type clientBuilder struct {
	r relay.Builder
}

func NewClientBuilder() ClientBuilder {
	return &clientBuilder{r: relay.NewBuilder()}
}

func (b *clientBuilder) Relay() relay.Builder {
	return b.r
}

func (b *clientBuilder) Build(close func()) Client {
	return NewClient(b.r.Build(), close)
}

type client struct {
	r     relay.Relay
	close func()
}

func NewClient(r relay.Relay, close func()) Client {
	return &client{r: r, close: close}
}

func (r *client) Relay() relay.Relay {
	return r.r
}

func (r *client) Close() {
	r.close()
}
