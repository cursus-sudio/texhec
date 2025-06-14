package src

import (
	"backend/src/backendapi"
	"backend/src/modules"
	"backend/src/utils"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	utilsPkg utils.Pkg,
	modulesPkg modules.Pkg,
	backendApiPkg backendapi.Pkg,
) ioc.Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			utilsPkg,
			modulesPkg,
			backendApiPkg,
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
