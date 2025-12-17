package test

import (
	"engine/modules/hierarchy"
	"engine/modules/hierarchy/pkg"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	World ecs.World
	Tool  hierarchy.Interface
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		hierarchypkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	world := ecs.NewWorld()
	toolFactory := ioc.Get[ecs.ToolFactory[hierarchy.HierarchyTool]](c)
	tool := toolFactory.Build(world)
	return Setup{
		world,
		tool.Hierarchy(),
	}
}
