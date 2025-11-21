package test

import (
	"frontend/modules/transform"
	transformpkg "frontend/modules/transform/pkg"
	"shared/services/clock"
	"shared/services/ecs"
	"shared/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	World       ecs.World
	Tool        transform.TransformTool
	Transaction transform.TransformTransaction
}

func NewSetup() Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		transformpkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	world := ecs.NewWorld()
	toolFactory := ioc.Get[ecs.ToolFactory[transform.TransformTool]](c)
	tool := toolFactory.Build(world)
	return Setup{
		world,
		tool,
		tool.Transaction(),
	}
}
