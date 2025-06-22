package backendconnector

import (
	"backend/services/api"
	"frontend/services/backendconnector/localconnector"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs        []ioc.Pkg
	loadDefault func(c ioc.Dic) api.Server
}

func Package(
	localPkg localconnector.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			localPkg,
		},
		loadDefault: func(c ioc.Dic) api.Server {
			return ioc.Get[localconnector.Connector](c).Connect()
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Backend {
		b := NewBuilder()
		b.DefaultConnection(func() api.Server { return pkg.loadDefault(c) })
		return b.Build()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) api.Server {
		return ioc.Get[Backend](c).Connection()
	})
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
