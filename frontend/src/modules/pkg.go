package modules

import (
	"frontend/src/modules/backendconnector"

	"github.com/ogiusek/ioc"
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

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
