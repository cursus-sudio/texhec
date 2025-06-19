package localconnector

import (
	backendsrc "backend/src"
	"backend/src/backendapi"

	"github.com/ogiusek/ioc/v2"
)

type Connector interface {
	Connect() backendapi.Backend
}

type connector struct {
	backend backendapi.Backend
}

func newConnector(backend backendapi.Backend) Connector {
	return &connector{
		backend: backend,
	}
}

func (connector *connector) Connect() backendapi.Backend {
	return connector.backend
}

type Pkg struct {
	backend backendapi.Backend
}

func Package(
	backend backendsrc.Pkg,
) Pkg {
	b := ioc.NewBuilder()
	backend.Register(b)
	c := b.Build()
	return Pkg{
		backend: ioc.Get[backendapi.Backend](c),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Connector {
		return newConnector(pkg.backend)
	})
}
