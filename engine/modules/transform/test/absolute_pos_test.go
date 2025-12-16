package test

import (
	"engine/modules/transform"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestAbsolutePos(t *testing.T) {
	setup := NewSetup(t)
	entity := setup.World.NewEntity()

	setup.Transform.Pos().SaveComponent(entity, transform.NewPos(10, 10, 10))
	setup.expectAbsolutePos(entity, transform.NewPos(10, 10, 10))

	setup.Transform.Pos().SaveComponent(entity, transform.NewPos(15, 15, 15))
	setup.expectAbsolutePos(entity, transform.NewPos(15, 15, 15))

	setup.Transform.SetAbsolutePos(entity, transform.AbsolutePosComponent{Pos: mgl32.Vec3{5, 5, 5}})
	setup.expectAbsolutePos(entity, transform.NewPos(5, 5, 5))
}
