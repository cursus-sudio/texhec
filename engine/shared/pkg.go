package shared

import (
	"shared/services/api"
	"shared/services/clock"
	"shared/services/codec"
	"shared/services/events"
	"shared/services/logger"
	"shared/services/runtime"
	"shared/services/uuid"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	apiPkg api.Pkg,
	clockPkg clock.Pkg,
	loggerPkg logger.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			apiPkg,
			clockPkg,
			events.Package(),
			loggerPkg,
			codec.Package(),
			runtime.Package(),
			uuid.Package(),
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
