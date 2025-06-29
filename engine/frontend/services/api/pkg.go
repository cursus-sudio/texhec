package api

import (
	"frontend/services/api/tcp"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(tcpPkg tcp.Pkg) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			tcpPkg,
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
	ioc.RegisterSingleton(b, func(c ioc.Dic) connection.Definition {
		return connection.NewDefinition()
	})

	ioc.RegisterTransient(b, func(c ioc.Dic) connection.Connection {
		return ioc.Get[connection.Definition](c).Build()
	})
}
