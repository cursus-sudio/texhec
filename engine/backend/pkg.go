package backend

import (
	"backend/services/api"
	"backend/services/clock"
	"backend/services/db"
	"backend/services/events"
	"backend/services/files"
	"backend/services/logger"
	"backend/services/saves"
	"backend/services/scopes"
	"backend/services/uuid"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	clockPkg clock.Pkg,
	dbPkg db.Pkg,
	filesPkg files.Pkg,
	modsPkgs []ioc.Pkg,
) Pkg {
	return Pkg{
		pkgs: append([]ioc.Pkg{
			api.Package(),
			clockPkg,
			dbPkg,
			events.Package(),
			filesPkg,
			logger.Package(),
			saves.Package(),
			scopes.Package(),
			uuid.Package(),
		}, modsPkgs...),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
