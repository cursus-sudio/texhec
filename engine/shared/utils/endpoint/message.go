package endpoint

import (
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type Message interface {
	relay.Message
	C() ioc.Dic
	UseC(ioc.Dic)
}

type message struct {
	relay.Message
	c ioc.Dic
}

func (mess *message) C() ioc.Dic     { return mess.c }
func (mess *message) UseC(c ioc.Dic) { mess.c = c }
func NewMessage() Message            { return &message{} }

type messageEndpoint[Mess Message] interface {
	Handle(Mess)
}

func MessageRegister[Endpoint messageEndpoint[Req], Req Message](b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, s connection.Definition) connection.Definition {
		relay.MessageRegister(s.MessageListener().Relay(), func(req Req) {
			c := req.C()
			endpoint := ioc.GetServices[Endpoint](c)
			endpoint.Handle(req)
		})
		return s
	})
}
