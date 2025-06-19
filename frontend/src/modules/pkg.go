package modules

import (
	"frontend/src/modules/backendconnector"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	backendPkg backendconnector.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			backendPkg,
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
