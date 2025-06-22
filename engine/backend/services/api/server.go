package api

import "github.com/ogiusek/relay/v2"

// services used by client

type Server interface { // client -> server
	Relay() relay.Relay
	Close()
}

type ServerBuilder interface { // client <- server
	Relay() relay.Builder
	Build(close func()) Server
}

type serverBuilder struct {
	r relay.Builder
}

func NewServerBuilder() ServerBuilder {
	return &serverBuilder{r: relay.NewBuilder()}
}

func (b *serverBuilder) Relay() relay.Builder {
	return b.r
}

func (b *serverBuilder) Build(close func()) Server {
	return NewServer(b.r.Build(), close)
}

type server struct {
	r     relay.Relay
	close func()
}

func NewServer(r relay.Relay, close func()) Server {
	return &server{r: r, close: close}
}

func (r *server) Relay() relay.Relay {
	return r.r
}

func (r *server) Close() {
	r.close()
}
