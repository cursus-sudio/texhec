package frontend

import (
	"frontend/services/backendconnector"
	"frontend/services/console"
	"frontend/services/draw"
	"frontend/services/ecs"
	"frontend/services/events"
	"frontend/services/inputs"
	"frontend/services/scenes"
	"frontend/services/window"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	backendConnectorPkg backendconnector.Pkg,
	inputsPkg inputs.Pkg,
	windowPkg window.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			backendConnectorPkg,
			inputsPkg,
			windowPkg,
			draw.Package(),
			console.Package(),
			ecs.Package(),
			events.Package(),
			scenes.Package(),
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}

// type Pkg struct {
// 	pkgs []ioc.Pkg
// }
//
// func Package() Pkg {
// 	// var backendPkg backend.Pkg = backend.Package(
// 	// 	backendapi.Package(),
// 	// 	clock.Package(time.RFC3339Nano),
// 	// 	db.Package(
// 	// 		fmt.Sprintf("%s/db.sql", userStorage),
// 	// 		fmt.Sprintf("%s/engine/backend/db/migrations", currentDir),
// 	// 	),
// 	// 	files.Package(fmt.Sprintf("%s/files", userStorage)),
// 	// 	logger.Package(),
// 	// 	saves.Package(),
// 	// 	scopecleanup.Package(),
// 	// 	uuid.Package(),
// 	// 	[]ioc.Pkg{
// 	// 		exBackendModPkg{},
// 	// 		Package(),
// 	// 	},
// 	// )
// 	//
// 	// var pkg frontendsrc.Pkg = frontendsrc.Package(
// 	// 	services.Package(
// 	// 		backendconnector.Package(localconnector.Package(backendPkg)),
// 	// 		inputs.Package(),
// 	// 		window.Package(),
// 	// 	),
// 	// )
// 	return Pkg{
// 		pkgs: []ioc.Pkg{},
// 	}
// }
//
// func (pkg Pkg) Register(b ioc.Builder) {
// 	for _, pkg := range pkg.pkgs {
// 		pkg.Register(b)
// 	}
// }
