package backendconnection

import (
	"frontend/services/backendconnection/localconnector"
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs        []ioc.Pkg
	loadDefault func(c ioc.Dic) connection.Connection
}

func Package(
	localPkg localconnector.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			localPkg,
		},
		loadDefault: func(c ioc.Dic) connection.Connection {
			return ioc.Get[localconnector.Connector](c).Connect()
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Backend {
		b := NewBuilder()
		b.DefaultConnection(func() connection.Connection { return pkg.loadDefault(c) })
		return b.Build()
	})

	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
