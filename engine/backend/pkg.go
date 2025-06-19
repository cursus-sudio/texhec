package backend

import (
	"backend/services/backendapi"
	"backend/services/clock"
	"backend/services/db"
	"backend/services/files"
	"backend/services/logger"
	"backend/services/saves"
	"backend/services/scopecleanup"
	"backend/services/uuid"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	backendApiPkg backendapi.Pkg,
	clockPkg clock.Pkg,
	dbPkg db.Pkg,
	filesPkg files.Pkg,
	loggerPkg logger.Pkg,
	savesPkg saves.Pkg,
	scopeCleanUpPkg scopecleanup.Pkg,
	uuidPkg uuid.Pkg,
	modsPkgs []ioc.Pkg,
) Pkg {
	return Pkg{
		pkgs: append([]ioc.Pkg{
			backendApiPkg,
			clockPkg,
			dbPkg,
			filesPkg,
			loggerPkg,
			savesPkg,
			scopeCleanUpPkg,
			uuidPkg,
		}, modsPkgs...),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
