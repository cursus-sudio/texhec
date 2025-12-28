package test

import (
	"engine/modules/transform"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestHierarchy(t *testing.T) {
	setup := NewSetup(t)
	parent := setup.NewEntity()
	child := setup.NewEntity()

	setup.Hierarchy().SetParent(child, parent)
	setup.Transform().Pos().Set(parent, transform.NewPos(5, 5, 5))
	setup.Transform().Parent().Set(child, transform.NewParent(transform.RelativePos))

	setup.Transform().Pos().Set(child, transform.NewPos(5, 5, 5))
	setup.expectAbsolutePos(child, transform.NewPos(10, 10, 10))

	setup.Transform().AbsolutePos().Set(child, transform.AbsolutePosComponent{Pos: mgl32.Vec3{5, 5, 5}})
	setup.expectAbsolutePos(child, transform.NewPos(5, 5, 5))

	setup.Transform().Pos().Set(child, transform.NewPos(10, 10, 10))
	setup.expectAbsolutePos(child, transform.NewPos(15, 15, 15))

	setup.Transform().Pos().Set(parent, transform.NewPos(10, 10, 10))
	setup.expectAbsolutePos(child, transform.NewPos(20, 20, 20))

	setup.Transform().Pos().Set(child, transform.NewPos(0, 0, 0))
	setup.expectAbsolutePos(child, transform.NewPos(10, 10, 10))
}
