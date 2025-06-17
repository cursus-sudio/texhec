package localconnector

import (
	backendsrc "backend/src"
	"backend/src/backendapi"

	"github.com/ogiusek/ioc"
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
	c := ioc.NewContainer()
	backend.Register(c)
	return Pkg{
		backend: ioc.Get[backendapi.Backend](c),
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) Connector {
		return newConnector(pkg.backend)
	})
}
