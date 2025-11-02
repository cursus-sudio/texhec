package api

import (
	"backend/services/scopes"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterScoped(b, scopes.UserSession, func(c ioc.Dic) connection.Definition {
		return connection.NewDefinition()
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) connection.Connection {
		return ioc.Get[connection.Definition](c).Build()
	})
	ioc.RegisterDependency[connection.Connection, connection.Definition](b)
}
