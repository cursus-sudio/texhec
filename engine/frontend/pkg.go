package frontend

import (
	"frontend/example/ping"
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
			// TODO TEMP
			ping.ClientPackage(),
			//
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
