package utils

import (
	"backend/src/utils/clock"
	"backend/src/utils/db"
	"backend/src/utils/files"
	"backend/src/utils/logger"
	"backend/src/utils/services"
	"backend/src/utils/uuid"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	clockPkg clock.Pkg,
	dbPkg db.Pkg,
	filesPkg files.Pkg,
	loggerPkg logger.Pkg,
	servicesPkg services.Pkg,
	uuidPkg uuid.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			clockPkg,
			dbPkg,
			filesPkg,
			loggerPkg,
			servicesPkg,
			uuidPkg,
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
