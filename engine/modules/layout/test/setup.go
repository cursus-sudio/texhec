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
	World     ecs.World         `inject:"1"`
	Hierarchy hierarchy.Service `inject:"1"`
	Transform transform.Service `inject:"1"`
	Layout    layout.Service    `inject:"1"`
	T         *testing.T
}

func NewSetup(t *testing.T) Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		ecs.Package(),
		hierarchypkg.Package(),
		transformpkg.Package(),
		layoutpkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	setup := ioc.GetServices[Setup](c)
	setup.T = t
	return setup
}

func (s Setup) Expect(entity ecs.EntityID, x, y float32) {
	s.T.Helper()
	expected := mgl32.Vec3{x, y, 1}
	pos, _ := s.Transform.AbsolutePos().Get(entity)
	if pos.Pos != expected {
		pivot, _ := s.Transform.PivotPoint().Get(entity)
		parentPivot, _ := s.Transform.ParentPivotPoint().Get(entity)
		size, _ := s.Transform.AbsoluteSize().Get(entity)

		parent, _ := s.Hierarchy.Parent(entity)
		pSize, _ := s.Transform.AbsoluteSize().Get(parent)
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
