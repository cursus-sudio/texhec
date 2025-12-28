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
	"testing"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	ecs.World
	hierarchy.HierarchyTool
	transform.TransformTool
	layout.LayoutTool
	T *testing.T
}

func NewSetup(t *testing.T) Setup {
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
		T:     t,
	}
	setup.HierarchyTool = ioc.Get[hierarchy.ToolFactory](c).Build(setup)
	setup.TransformTool = ioc.Get[transform.ToolFactory](c).Build(setup)
	setup.LayoutTool = ioc.Get[layout.ToolFactory](c).Build(setup)
	return setup
}

func (s Setup) Expect(entity ecs.EntityID, x, y float32) {
	s.T.Helper()
	expected := mgl32.Vec3{x, y}
	pos, _ := s.Transform().AbsolutePos().Get(entity)
	if pos.Pos != expected {
		pivot, _ := s.Transform().PivotPoint().Get(entity)
		parentPivot, _ := s.Transform().ParentPivotPoint().Get(entity)
		size, _ := s.Transform().AbsoluteSize().Get(entity)

		parent, _ := s.Hierarchy().Parent(entity)
		pSize, _ := s.Transform().AbsoluteSize().Get(parent)
		s.T.Errorf(
			"expected %v and got %v (pivot %v, parent %v, size %v, pSize %v)",
			expected,
			pos,
			pivot,
			parentPivot,
			size,
			pSize,
		)
	}
}
