package backendapipkg

import (
	"backend/src/backendapi"
	"backend/src/backendapi/ping"
	"backend/src/backendapi/tacticalmapapi"
	"backend/src/utils/httperrors"

	"github.com/ogiusek/ioc/v2"
	"github.com/ogiusek/relay/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package() Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			ping.Pkg{},
			tacticalmapapi.Pkg{},
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
	ioc.RegisterSingleton(b, func(c ioc.Dic) backendapi.Builder {
		b := backendapi.NewBuilder()
		r := b.Relay()
		r.DefaultHandler(func(req any) (relay.Res, error) { return nil, httperrors.Err404 })
		return b
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) backendapi.Backend {
		return ioc.Get[backendapi.Builder](c).Build()
	})
}
