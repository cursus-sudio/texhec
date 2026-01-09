package test

import (
	"engine/modules/hierarchy"
	hierarchypkg "engine/modules/hierarchy/pkg"
	"engine/modules/transform"
	transformpkg "engine/modules/transform/pkg"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"testing"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type world struct {
	ecs.World
	hierarchy.HierarchyTool
	transform.TransformTool
}

type Setup struct {
	world
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
	world := world{
		World: ecs.NewWorld(),
	}
	world.HierarchyTool = ioc.Get[hierarchy.ToolFactory](c).Build(world)
	world.TransformTool = ioc.Get[transform.ToolFactory](c).Build(world)
	return Setup{world}
}

func (setup Setup) expectAbsolutePos(t *testing.T, entity ecs.EntityID, expectedPos transform.PosComponent) {
	t.Helper()
	pos, _ := setup.Transform().AbsolutePos().Get(entity)
	if pos.Pos != expectedPos.Pos {
		t.Errorf("expected pos %v but has %v", expectedPos, pos)
	}
}

func (setup Setup) expectAbsoluteSize(t *testing.T, entity ecs.EntityID, expectedSize transform.SizeComponent) {
	t.Helper()
	size, _ := setup.Transform().AbsoluteSize().Get(entity)
	if size.Size != expectedSize.Size {
		t.Errorf("expected size %v but has %v", expectedSize, size)
	}
}
