package test

import (
	"engine/modules/hierarchy"
	hierarchypkg "engine/modules/hierarchy/pkg"
	"engine/modules/transform"
	transformpkg "engine/modules/transform/pkg"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	World          ecs.World
	Tool           transform.Tool
	Transaction    transform.Transaction
	HierarchyArray ecs.ComponentsArray[hierarchy.ParentComponent]
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		hierarchypkg.Package(),
		transformpkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	world := ecs.NewWorld()
	toolFactory := ioc.Get[ecs.ToolFactory[transform.Tool]](c)
	tool := toolFactory.Build(world)
	hierarchyArray := ecs.GetComponentsArray[hierarchy.ParentComponent](world)
	return Setup{
		world,
		tool,
		tool.Transaction(),
		hierarchyArray,
	}
}
