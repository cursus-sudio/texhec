package test

import (
	"testing"
)

func BenchmarkAddNChildrenWithParent(b *testing.B) {
	setup := NewSetup()
	grandParent := setup.World.NewEntity()
	parent := grandParent
	parentCount := 0
	for i := 0; i < parentCount; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)
		parent = child
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)
	}

	setup.Tool.FlatChildren(parent)
}

func BenchmarkAddNChildrenWith5Parents(b *testing.B) {
	setup := NewSetup()
	grandParent := setup.World.NewEntity()
	parent := grandParent
	parentCount := 5
	for i := 0; i < parentCount; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)
		parent = child
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)
	}

	setup.Tool.FlatChildren(parent)
}
