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
	World     ecs.World
	Transform transform.Interface
	Hierarchy hierarchy.Interface
}

type Setup struct {
	world
	T *testing.T
}

func NewSetup(t *testing.T) Setup {
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
	w := ecs.NewWorld()
	return Setup{
		world{
			w,
			ioc.Get[ecs.ToolFactory[transform.TransformTool]](c).Build(w).Transform(),
			ioc.Get[ecs.ToolFactory[hierarchy.HierarchyTool]](c).Build(w).Hierarchy(),
		},
		t,
	}
}

func (setup Setup) expectAbsolutePos(entity ecs.EntityID, expectedPos transform.PosComponent) {
	pos, _ := setup.Transform.AbsolutePos().Get(entity)
	if pos.Pos != expectedPos.Pos {
		setup.T.Errorf("expected pos %v but has %v", expectedPos, pos)
	}
}

func (setup Setup) expectAbsoluteSize(entity ecs.EntityID, expectedSize transform.SizeComponent) {
	size, _ := setup.Transform.AbsoluteSize().Get(entity)
	if size.Size != expectedSize.Size {
		setup.T.Errorf("expected size %v but has %v", expectedSize, size)
	}
}
