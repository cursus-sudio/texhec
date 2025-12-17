package test

import (
	"engine/modules/transform"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestSize(t *testing.T) {
	setup := NewSetup(t)
	entity := setup.World.NewEntity()

	setup.Transform.Size().Set(entity, transform.NewSize(10, 10, 10))
	setup.expectAbsoluteSize(entity, transform.NewSize(10, 10, 10))

	setup.Transform.Size().Set(entity, transform.NewSize(15, 15, 15))
	setup.expectAbsoluteSize(entity, transform.NewSize(15, 15, 15))

	setup.Transform.SetAbsoluteSize(entity, transform.AbsoluteSizeComponent{Size: mgl32.Vec3{5, 5, 5}})
	setup.expectAbsoluteSize(entity, transform.NewSize(5, 5, 5))
}
