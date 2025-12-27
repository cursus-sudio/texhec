package test

import (
	"engine/modules/hierarchy"
	hierarchypkg "engine/modules/hierarchy/pkg"
	"engine/modules/layout"
	layoutpkg "engine/modules/layout/pkg"
	"engine/modules/transform"
	transformpkg "engine/modules/transform/pkg"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	ecs.World
	hierarchy.HierarchyTool
	transform.TransformTool
	layout.LayoutTool
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		hierarchypkg.Package(),
		transformpkg.Package(),
		layoutpkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	setup := Setup{
		World: ecs.NewWorld(),
	}
	setup.HierarchyTool = ioc.Get[hierarchy.ToolFactory](c).Build(setup)
	setup.TransformTool = ioc.Get[transform.ToolFactory](c).Build(setup)
	setup.LayoutTool = ioc.Get[layout.ToolFactory](c).Build(setup)
	return setup
}
