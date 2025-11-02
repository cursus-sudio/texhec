package dbpkg

import (
	"embed"
	"frontend/services/scopes"
	"shared/services/db"

	"github.com/ogiusek/ioc/v2"
)

//go:embed migrations/*.sql
var migrations embed.FS

type pkg struct {
	dbPath string
	dbPkg  ioc.Pkg
}

func Package(dbPath string) ioc.Pkg {
	return pkg{
		dbPath: dbPath,
		dbPkg: db.Package(dbPath, migrations, scopes.Request, func(c ioc.Dic, cleanUp func(err error)) {
			s := ioc.Get[scopes.RequestService](c)
			s.AddCleanListener(func(args scopes.RequestEndArgs) {
				cleanUp(args.Error)
			})
		}),
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	pkg.dbPkg.Register(b)
	ioc.RegisterDependency[db.Tx, scopes.RequestService](b)
}
