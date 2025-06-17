package backendapi

import (
	"backend/src/backendapi/ping"
	"backend/src/backendapi/tacticalmapapi"
	"backend/src/utils/httperrors"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/relay"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

type RelayPkg interface {
	Register(relay.Relay)
}

func (pkg Pkg) Register(c ioc.Dic) {
	relayPackages := []RelayPkg{
		ping.Package(c),
		tacticalmapapi.Package(c),
	}

	ioc.RegisterSingleton(c, func(c ioc.Dic) Backend {
		r := relay.NewRelay(relay.NewConfigBuilder().
			DefaultHandler(func(req any) (relay.Res, error) { return nil, httperrors.Err404 }).
			Build())

		for _, pkg := range relayPackages {
			pkg.Register(r)
		}

		return newBackend(r)
	})
}
