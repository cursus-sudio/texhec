package connection

import "github.com/ogiusek/ioc/v2"

const (
	OrderEndpoint ioc.Order = iota
	OrderAuthorize
	OrderAuthenticate
	OrderAttachServices
)

type Definition interface {
	Close() CloseBuilder
	MessageListener() MessageListenerBuilder
	Build() Connection
}

type connectionBuilder struct {
	CloseBuilder
	MessageListenerBuilder
}

func NewDefinition() Definition {
	return connectionBuilder{
		CloseBuilder:           NewCloseBuilder(),
		MessageListenerBuilder: NewMessageListenerBuilder(),
	}
}

func (b connectionBuilder) Close() CloseBuilder                     { return b.CloseBuilder }
func (b connectionBuilder) MessageListener() MessageListenerBuilder { return b.MessageListenerBuilder }
func (b connectionBuilder) Build() Connection {
	return NewConnection(b.Close().Build(), b.MessageListener().Build())
}

//

type Connection interface {
	Closer
	MessagerEmitter
}

type connection struct {
	Closer
	MessagerEmitter
}

func NewConnection(closer Closer, messageEmitter MessagerEmitter) Connection {
	return &connection{
		Closer:          closer,
		MessagerEmitter: messageEmitter,
	}
}
