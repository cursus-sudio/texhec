package dbpkg

import (
	"embed"
	"frontend/services/scopes"
	"shared/services/db"

	"github.com/ogiusek/ioc/v2"
)

// type Pkg struct{}
//
// func Package() Pkg {
// 	return Pkg{}
// }
//
// func (Pkg) Register(c ioc.Dic) {
//
// }

//go:embed migrations/*.sql
var migrations embed.FS

type Pkg struct {
	dbPath string
	dbPkg  db.Pkg
}

func Package(dbPath string) Pkg {
	return Pkg{
		dbPath: dbPath,
		dbPkg: db.Package(dbPath, migrations, scopes.Request, func(c ioc.Dic, cleanUp func(err error)) {
			s := ioc.Get[scopes.RequestService](c)
			s.AddCleanListener(func(args scopes.RequestEndArgs) {
				cleanUp(args.Error)
			})
		}),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	pkg.dbPkg.Register(b)
	ioc.RegisterDependency[db.Tx, scopes.RequestService](b)
}
