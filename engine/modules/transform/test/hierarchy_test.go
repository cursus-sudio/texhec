package test

import (
	"engine/modules/transform"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestHierarchy(t *testing.T) {
	setup := NewSetup()
	parent := setup.NewEntity()
	child := setup.NewEntity()

	setup.hierarchy.SetParent(child, parent)
	setup.transform.Pos().Set(parent, transform.NewPos(5, 5, 5))
	setup.transform.Parent().Set(child, transform.NewParent(transform.RelativePos))

	setup.transform.Pos().Set(child, transform.NewPos(5, 5, 5))
	setup.expectAbsolutePos(t, child, transform.NewPos(10, 10, 10))

	setup.transform.AbsolutePos().Set(child, transform.AbsolutePosComponent{Pos: mgl32.Vec3{5, 5, 5}})
	setup.expectAbsolutePos(t, child, transform.NewPos(5, 5, 5))

	setup.transform.Pos().Set(child, transform.NewPos(10, 10, 10))
	setup.expectAbsolutePos(t, child, transform.NewPos(15, 15, 15))

	setup.transform.Pos().Set(parent, transform.NewPos(10, 10, 10))
	setup.expectAbsolutePos(t, child, transform.NewPos(20, 20, 20))

	setup.transform.Pos().Set(child, transform.NewPos(0, 0, 0))
	setup.expectAbsolutePos(t, child, transform.NewPos(10, 10, 10))
}
