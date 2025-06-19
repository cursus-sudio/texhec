package backendapi

import (
	"backend/utils/httperrors"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		b := NewBuilder()
		r := b.Relay()
		r.DefaultHandler(func(req any) (relay.Res, error) { return nil, httperrors.Err404 })
		return b
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) Backend {
		return ioc.Get[Builder](c).Build()
	})
}
