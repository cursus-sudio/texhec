package backend

import (
	"backend/services/api"
	"backend/services/clients"
	"backend/services/db"
	"backend/services/files"
	"backend/services/logger"
	"backend/services/saves"
	"backend/services/scopes"
	"shared"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	sharedPkg shared.Pkg,
	dbPkg db.Pkg,
	filesPkg files.Pkg,
	loggerPkg logger.Pkg,
	modsPkgs []ioc.Pkg,
) Pkg {
	return Pkg{
		pkgs: append([]ioc.Pkg{
			sharedPkg,
			api.Package(),
			clients.Package(),
			dbPkg,
			filesPkg,
			loggerPkg,
			saves.Package(),
			scopes.Package(),
		}, modsPkgs...),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
