package test

import (
	"engine/modules/transform"
	"engine/services/ecs"
	"testing"
)

func BenchmarkSetPos(b *testing.B) {
	setup := NewSetup()

	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		pos := transform.NewPos(0, 0, float32(i))
		setup.Transform().AbsolutePos().Set(entity, transform.AbsolutePosComponent(pos))
	}
}

func BenchmarkGetPos(b *testing.B) {
	setup := NewSetup()
	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		setup.Transform().AbsolutePos().Get(entity)
	}
}

func BenchmarkTransformGetPos(b *testing.B) {
	setup := NewSetup()
	entity := setup.NewEntity()
	transform := setup.Transform()
	for i := 0; i < b.N; i++ {
		transform.AbsolutePos().Get(entity)
	}
}

func BenchmarkArrGetPos(b *testing.B) {
	setup := NewSetup()
	entity := setup.NewEntity()
	arr := setup.Transform().AbsolutePos()
	for i := 0; i < b.N; i++ {
		arr.Get(entity)
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

func BenchmarkSetAndGetPos(b *testing.B) {
	setup := NewSetup()

	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		pos := transform.NewPos(0, 0, float32(i))
		setup.Transform().AbsolutePos().Set(entity, transform.AbsolutePosComponent(pos))
		for i := 0; i < 1; i++ {
			setup.Transform().AbsolutePos().Get(entity)
		}
	}
}

func BenchmarkSetAndDoubleGetPos(b *testing.B) {
	setup := NewSetup()

	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		pos := transform.NewPos(0, 0, float32(i))
		setup.Transform().AbsolutePos().Set(entity, transform.AbsolutePosComponent(pos))
		for i := 0; i < 2; i++ {
			setup.Transform().AbsolutePos().Get(entity)
		}
	}
}

func BenchmarkSetAndTripleGetPos(b *testing.B) {
	setup := NewSetup()

	entity := setup.NewEntity()
	for i := 0; i < b.N; i++ {
		pos := transform.NewPos(0, 0, float32(i))
		setup.Transform().AbsolutePos().Set(entity, transform.AbsolutePosComponent(pos))
		for i := 0; i < 3; i++ {
			setup.Transform().AbsolutePos().Get(entity)
		}
	}
}
