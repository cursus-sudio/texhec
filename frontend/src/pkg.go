package src

import (
	"frontend/src/engine"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	enginePkg engine.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			enginePkg,
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
