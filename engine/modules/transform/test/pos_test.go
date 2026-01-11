package test

import (
	"engine/modules/transform"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestPos(t *testing.T) {
	setup := NewSetup()
	entity := setup.NewEntity()

	setup.transform.Pos().Set(entity, transform.NewPos(10, 10, 10))
	setup.expectAbsolutePos(t, entity, transform.NewPos(10, 10, 10))

	setup.transform.Pos().Set(entity, transform.NewPos(15, 15, 15))
	setup.expectAbsolutePos(t, entity, transform.NewPos(15, 15, 15))

	setup.transform.AbsolutePos().Set(entity, transform.AbsolutePosComponent{Pos: mgl32.Vec3{5, 5, 5}})
	setup.expectAbsolutePos(t, entity, transform.NewPos(5, 5, 5))
}
