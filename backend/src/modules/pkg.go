package modules

import (
	"backend/src/modules/saves"
	"backend/src/modules/tacticalmap"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	savesPackage saves.Pkg,
	tacticalMapPackage tacticalmap.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			savesPackage,
			tacticalMapPackage,
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
