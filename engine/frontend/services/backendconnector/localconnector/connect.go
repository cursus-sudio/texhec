package localconnector

import (
	"backend"
	"backend/services/backendapi"

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
	backend backend.Pkg,
) Pkg {
	cBuilder := ioc.NewBuilder()
	backend.Register(cBuilder)
	c := cBuilder.Build()
	backendApi := ioc.Get[backendapi.Backend](c)
	return Pkg{
		backend: backendApi,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Connector {
		return newConnector(pkg.backend)
	})
}
