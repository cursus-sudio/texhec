package backendconnector

import (
	"backend/src/backendapi"
	"frontend/src/modules/backendconnector/localconnector"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs        []ioc.Pkg
	loadDefault func(c ioc.Dic)
}

func Package(
	localPkg localconnector.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			localPkg,
		},
		loadDefault: func(c ioc.Dic) {
			ioc.Get[localconnector.Connector](c).Connect()
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) Backend {
		return NewBackend(func() backendapi.Backend {
			return ioc.Get[localconnector.Connector](c).Connect()
		})
	})
	ioc.RegisterSingleton(c, func(c ioc.Dic) backendapi.Backend {
		return ioc.Get[Backend](c).Connection()
	})
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
	pkg.loadDefault(c)
}
