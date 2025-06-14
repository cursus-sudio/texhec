package engine

import (
	"frontend/src/engine/console"
	"frontend/src/engine/draw"
	"frontend/src/engine/ecs"
	"frontend/src/engine/inputs"
	"frontend/src/engine/scenes"
	"frontend/src/engine/window"

	"github.com/ogiusek/ioc"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	inputsPkg inputs.Pkg,
	windowPkg window.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			inputsPkg,
			windowPkg,
			draw.Package(),
			console.Package(),

			ecs.Package(),
			scenes.Package(),
		},
	}
}

func (pkg Pkg) Register(c ioc.Dic) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(c)
	}
}
