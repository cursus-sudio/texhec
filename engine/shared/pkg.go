package shared

import (
	"shared/services/clock"
	"shared/services/events"
	"shared/services/requestcodec"
	"shared/services/uuid"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	clockPkg clock.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			uuid.Package(),
			events.Package(),
			requestcodec.Package(),
			clockPkg,
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
