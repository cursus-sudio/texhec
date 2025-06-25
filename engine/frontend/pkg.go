package frontend

import (
	"frontend/services/api"
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/draw"
	"frontend/services/ecs"
	"frontend/services/inputs"
	"frontend/services/scenes"
	"frontend/services/window"
	"shared"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	sharedPkg shared.Pkg,
	backendConnectorPkg backendconnection.Pkg,
	inputsPkg inputs.Pkg,
	windowPkg window.Pkg,
	mods []ioc.Pkg,
) Pkg {
	return Pkg{
		pkgs: append([]ioc.Pkg{
			sharedPkg,
			api.Package(),
			backendConnectorPkg,
			inputsPkg,
			windowPkg,
			draw.Package(),
			console.Package(),
			ecs.Package(),
			scenes.Package(),
		}, mods...),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
