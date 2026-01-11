package test

import (
	"engine/services/ecs"
	"testing"
)

func BenchmarkAddNChildrenWithParent(b *testing.B) {
	setup := NewSetup()
	grandParent := setup.World.NewEntity()
	parent := grandParent
	parentCount := 0
	for i := 0; i < parentCount; i++ {
		child := setup.World.NewEntity()
		setup.Service.SetParent(child, parent)
		parent = child
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child := setup.World.NewEntity()
		setup.Service.SetParent(child, parent)
	}
}

func BenchmarkAddNChildrenWith5Parents(b *testing.B) {
	setup := NewSetup()
	grandParent := setup.World.NewEntity()
	parent := grandParent
	parentCount := 5
	for i := 0; i < parentCount; i++ {
		child := setup.World.NewEntity()
		setup.Service.SetParent(child, parent)
		parent = child
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child := setup.World.NewEntity()
		setup.Service.SetParent(child, parent)
	}
}

func BenchmarkRemoveNChildren(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()

	children := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		children[i] = setup.World.NewEntity()
		setup.Service.SetParent(children[i], parent)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setup.World.RemoveEntity(children[i])
	}
}

func BenchmarkRemoveParentWithNChildren(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()

	for i := 0; i < b.N; i++ {
		child := setup.World.NewEntity()
		setup.Service.SetParent(child, parent)
	}

	b.ResetTimer()
	setup.World.RemoveEntity(parent)
}
