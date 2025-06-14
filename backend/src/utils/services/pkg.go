package services

import (
	"backend/src/utils/services/scopecleanup"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	scopeCleanUp scopecleanup.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			scopeCleanUp,
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
