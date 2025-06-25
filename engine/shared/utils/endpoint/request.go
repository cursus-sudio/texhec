package endpoint

import (
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type AnyRequest interface {
	C() ioc.Dic
	UseC(ioc.Dic)
}

type anyRequest struct{ c ioc.Dic }

func (req *anyRequest) C() ioc.Dic     { return req.c }
func (req *anyRequest) UseC(c ioc.Dic) { req.c = c }

func NewAnyRequest() AnyRequest { return &anyRequest{} }

type Request[Res relay.Res] interface {
	relay.Req[Res]
	AnyRequest
}

type request[Res relay.Res] struct {
	relay.Req[Res]
	AnyRequest
}

func NewRequest[Res relay.Res]() Request[Res] { return request[Res]{AnyRequest: NewAnyRequest()} }

type endpoint[Req Request[Res], Res relay.Res] interface {
	Handle(Req) (Res, error)
}

func Register[Endpoint endpoint[Req, Res], Req Request[Res], Res relay.Res](b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, s connection.Definition) connection.Definition {
		relay.Register(s.MessageListener().Relay(), func(req Req) (Res, error) {
			c := req.C()
			endpoint := ioc.GetServices[Endpoint](c)
			res, err := endpoint.Handle(req)
			return res, err
		})
		return s
	})
}
