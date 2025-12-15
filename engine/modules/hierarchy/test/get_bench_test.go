package test

import (
	"testing"
)

func BenchmarkChildren_1(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()

	setup.Tool.SetParent(child, parent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setup.Tool.Children(parent)
	}
}

func BenchmarkChildren_10(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()

	for i := 0; i < 10; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setup.Tool.Children(parent)
	}
}

func BenchmarkChildren_100(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()

	for i := 0; i < 100; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setup.Tool.Children(parent)
	}
}

func BenchmarkFlatChildren_1_1(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	setup.Tool.SetParent(child, parent)
	setup.Tool.SetParent(grandChild, child)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setup.Tool.FlatChildren(parent)
	}
}

func BenchmarkFlatChildren_10_10(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()

	for i := 0; i < 10; i++ {
		child := setup.World.NewEntity()
		setup.Tool.SetParent(child, parent)

		for j := 0; j < 10; j++ {
			grandChild := setup.World.NewEntity()
			setup.Tool.SetParent(grandChild, child)
		}
	}

	for i := 0; i < b.N; i++ {
		setup.Tool.FlatChildren(parent)
	}
}
