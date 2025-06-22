package localconnector

import (
	"backend"
	"backend/services/api"

	"github.com/ogiusek/ioc/v2"
)

type Connector interface {
	Connect() api.Server
}

type connector struct {
	backend api.Server
}

func newConnector(backend api.Server) Connector {
	return &connector{
		backend: backend,
	}
}

func (connector *connector) Connect() api.Server {
	return connector.backend
}

type Pkg struct {
	server        api.Server
	clientBuilder api.ClientBuilder
}

func Package(
	backend backend.Pkg,
) Pkg {
	var clientBuilder api.ClientBuilder = api.NewClientBuilder()
	cBuilder := ioc.NewBuilder()
	backend.Register(cBuilder)
	ioc.RegisterSingleton(cBuilder, func(c ioc.Dic) api.Client {
		return clientBuilder.Build(func() { print("write me\n\n\n\n\n\n\n\n\n\n\n\n\n") })
	})
	c := cBuilder.Build()
	server := ioc.Get[api.Server](c)
	return Pkg{
		server:        server,
		clientBuilder: clientBuilder,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) api.ClientBuilder {
		return pkg.clientBuilder
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) Connector { return newConnector(pkg.server) })
}
