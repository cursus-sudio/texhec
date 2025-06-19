package src

import (
	"backend/src/backendapi/backendapipkg"
	"backend/src/modules"
	"backend/src/utils"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	utilsPkg utils.Pkg,
	modulesPkg modules.Pkg,
	backendApiPkg backendapipkg.Pkg,
	modsPkgs []ioc.Pkg,
) Pkg {
	return Pkg{
		pkgs: append([]ioc.Pkg{
			utilsPkg,
			modulesPkg,
			backendApiPkg,
		}, modsPkgs...),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
