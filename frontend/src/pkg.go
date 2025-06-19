package src

import (
	"frontend/src/engine"
	"frontend/src/modules"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	enginePkg engine.Pkg,
	modules modules.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			enginePkg,
			modules,
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
