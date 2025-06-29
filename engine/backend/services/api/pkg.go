package api

import (
	"backend/services/api/tcp"
	"backend/services/scopes"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	tcpPkg tcp.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			tcpPkg,
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, service := range pkg.pkgs {
		service.Register(b)
	}
	ioc.RegisterScoped(b, scopes.UserSession, func(c ioc.Dic) connection.Definition {
		return connection.NewDefinition()
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) connection.Connection {
		return ioc.Get[connection.Definition](c).Build()
	})
}
