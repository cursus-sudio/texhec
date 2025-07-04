package backendconnection

import (
	"shared/utils/connection"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	loadDefault func(c ioc.Dic) connection.Connection
}

func Package(
	loadDefaults func(c ioc.Dic) connection.Connection,
) Pkg {
	return Pkg{
		loadDefault: loadDefaults,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Backend {
		b := NewBuilder()
		b.DefaultConnection(func() connection.Connection { return pkg.loadDefault(c) })
		return b.Build()
	})
}
