package clients

import (
	"backend/services/scopes"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(c ioc.Builder) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) ClientsBuilder {
		return NewClientsBuilder()
	})
	ioc.RegisterSingleton(c, func(c ioc.Dic) Clients {
		return ioc.Get[ClientsBuilder](c).Build()
	})
	ioc.RegisterScoped(c, scopes.UserSession, func(c ioc.Dic) SessionClient {
		return NewSessionClient(ioc.Get[Clients](c))
	})
}
