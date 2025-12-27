package test

import (
	"engine/modules/transform"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestLayout(t *testing.T) {
	setup := NewSetup()
	parent := setup.NewEntity()

	setup.Transform().Pos()

	setup.Transform().Pos().Set(parent, transform.NewPos(10, 10, 10))

	setup.Transform().Pos().Set(parent, transform.NewPos(15, 15, 15))

	setup.Transform().SetAbsolutePos(parent, transform.AbsolutePosComponent{Pos: mgl32.Vec3{5, 5, 5}})
}
