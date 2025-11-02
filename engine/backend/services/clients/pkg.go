package clients

import (
	"backend/services/scopes"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ClientsBuilder {
		return NewClientsBuilder()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) Clients {
		return ioc.Get[ClientsBuilder](c).Build()
	})
	ioc.RegisterDependency[Clients, ClientsBuilder](b)

	ioc.RegisterScoped(b, scopes.UserSession, func(c ioc.Dic) SessionClient {
		return NewSessionClient(ioc.Get[Clients](c))
	})
	ioc.RegisterDependency[SessionClient, Clients](b)
}
