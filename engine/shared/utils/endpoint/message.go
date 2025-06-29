package endpoint

import (
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type message interface {
	relay.Message
	anyRequest
}

type Message struct {
	relay.Message
	AnyRequest
}

func NewMessage() Message { return Message{} }

type messageEndpoint[Mess message] interface {
	Handle(Mess)
}

func MessageRegister[Endpoint messageEndpoint[Msg], Msg message](b ioc.Builder) {
	ioc.WrapService(b, connection.OrderEndpoint, func(c ioc.Dic, s connection.Definition) connection.Definition {
		relay.MessageRegister(s.MessageListener().Relay(), func(msg Msg) {
			c := msg.GetC()
			endpoint := ioc.GetServices[Endpoint](c)
			endpoint.Handle(msg)
		})
		return s
	})
}
