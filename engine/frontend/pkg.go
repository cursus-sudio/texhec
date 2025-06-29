package frontend

import (
	"frontend/services/api"
	"frontend/services/backendconnection"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media"
	"frontend/services/scenes"
	"frontend/services/scopes"
	"shared"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	sharedPkg shared.Pkg,
	apiPkg api.Pkg,
	backendConnectorPkg backendconnection.Pkg,
	mods []ioc.Pkg,
) Pkg {
	return Pkg{
		pkgs: append([]ioc.Pkg{
			sharedPkg,
			apiPkg,
			backendConnectorPkg,
			console.Package(),
			media.Package(), // media before ecs
			ecs.Package(),
			frames.Package(),
			scenes.Package(),
			scopes.Package(),
		}, mods...),
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
