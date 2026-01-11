package test

import (
	"engine/modules/transform"
	"engine/services/ecs"
	"testing"
)

func BenchmarkGetPos(b *testing.B) {
	setup := NewSetup()
	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		setup.transform.AbsolutePos().Get(entity)
	}
}

func BenchmarkRawGetPos(b *testing.B) {
	world := ecs.NewWorld()
	arr := ecs.GetComponentsArray[transform.AbsolutePosComponent](world) // no wrappers
	entity := world.NewEntity()
	for i := 0; i < b.N; i++ {
		arr.Get(entity)
	}
}

func BenchmarkSetAbsolutePos(b *testing.B) {
	setup := NewSetup()

	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		pos := transform.NewPos(0, 0, float32(i))
		setup.transform.AbsolutePos().Set(entity, transform.AbsolutePosComponent(pos))
	}
}

func BenchmarkSetAndGetAbsolutePos(b *testing.B) {
	setup := NewSetup()

	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		pos := transform.NewPos(0, 0, float32(i))
		setup.transform.Pos().Set(entity, pos)
		for i := 0; i < 1; i++ {
			setup.transform.AbsolutePos().Get(entity)
		}
	}
}
