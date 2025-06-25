package connection

import "github.com/ogiusek/relay/v2"

// builder

type CloseBuilder interface {
	OnClose(listener func())
	Build() Closer
}

type closeBuilder struct{ listeners []func() }

func NewCloseBuilder() CloseBuilder             { return &closeBuilder{listeners: []func(){}} }
func (b *closeBuilder) OnClose(listener func()) { b.listeners = append(b.listeners, listener) }
func (b *closeBuilder) Build() Closer {
	heapB := *b
	return &close{
		close: func() {
			for _, listener := range heapB.listeners {
				listener()
			}
		},
	}
}

// interface

type Closer interface {
	OnClose(func())
	Close()
}

type close struct {
	closed    bool
	listeners []func()
	close     func()
}

func (close *close) OnClose(listener func()) {
	close.listeners = append(close.listeners, listener)
}

func (close *close) Close() {
	if !close.closed {
		close.closed = true
		close.close()
		for _, listener := range close.listeners {
			listener()
		}
	}
}

// listener

// when we connect client to server we pass client Builder to server emmitter and vice versa
//
// when we want to introduce layer in between instead of passing other side builder we:
// 1. change builders to emit parsed format via way of communication
// 2. on reveived message we emit it via relay
//
// this way relay becomes ideal parser

type MessageListenerBuilder interface {
	Relay() relay.Builder
	Build() MessagerEmitter
}
type messageListenerBuilder struct {
	r relay.Builder
}

func NewMessageListenerBuilder() MessageListenerBuilder {
	return &messageListenerBuilder{r: relay.NewBuilder()}
}
func (b *messageListenerBuilder) Relay() relay.Builder   { return b.r }
func (b *messageListenerBuilder) Build() MessagerEmitter { return &messageEmmiter{r: b.r.Build()} }

//

type MessagerEmitter interface {
	Relay() relay.Relay
}
type messageEmmiter struct {
	r relay.Relay
}

func (emitter *messageEmmiter) Relay() relay.Relay { return emitter.r }
