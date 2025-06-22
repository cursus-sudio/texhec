package endpoint

import (
	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type AnyRequest interface {
	C() ioc.Dic
	UseC(ioc.Dic)
}

type Request[Res relay.Res] interface {
	relay.Req[Res]
	AnyRequest
}

type request[Res relay.Res] struct {
	relay.Req[Res] `json:"-"`
	c              ioc.Dic `json:"-"`
}

func NewRequest[Res relay.Res]() Request[Res] {
	return &request[Res]{}
}

func (req *request[Res]) C() ioc.Dic {
	return req.c
}

func (req *request[Res]) UseC(c ioc.Dic) {
	req.c = c
}

type endpoint[Req relay.Req[Res], Res relay.Res] interface {
	Handle(Req) (Res, error)
}

func Register[Endpoint endpoint[Req, Res], Req Request[Res], Res relay.Res](r relay.Builder) {
	relay.Register(r, func(req Req) (Res, error) {
		c := req.C()
		endpoint := ioc.GetServices[Endpoint](c)
		res, err := endpoint.Handle(req)
		return res, err
	})
}
