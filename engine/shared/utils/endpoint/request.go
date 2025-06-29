package endpoint

import (
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type request[Res relay.Res] interface {
	relay.Req[Res]
	anyRequest
}

type Request[Res relay.Res] struct {
	relay.Req[Res]
	AnyRequest
}

type endpoint[Req request[Res], Res relay.Res] interface {
	Handle(Req) (Res, error)
}

func Register[Endpoint endpoint[Req, Res], Req request[Res], Res relay.Res](b ioc.Builder) {
	ioc.WrapService(b, connection.OrderEndpoint, func(c ioc.Dic, s connection.Definition) connection.Definition {
		relay.Register(s.MessageListener().Relay(), func(req Req) (Res, error) {
			c := req.GetC()
			endpoint := ioc.GetServices[Endpoint](c)
			res, err := endpoint.Handle(req)
			return res, err
		})
		return s
	})
}
