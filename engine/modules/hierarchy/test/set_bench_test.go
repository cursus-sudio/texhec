package test

import (
	"testing"
)

func BenchmarkSetChildrenWithParent(b *testing.B) {
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

	parentChildren := setup.Tool.FlatChildren(parent)
	grandParentChildren := setup.Tool.FlatChildren(grandParent)
	parentLen := len(parentChildren.GetIndices())
	grandParentLen := len(grandParentChildren.GetIndices())
	if parentLen+parentCount != grandParentLen {
		b.Errorf(
			"flat children count of parent and grand parent doesn't match. expected %v and got %v",
			parentLen+parentCount,
			grandParentLen,
		)
	}
}

func BenchmarkSetChildrenWith5Parents(b *testing.B) {
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

	parentChildren := setup.Tool.FlatChildren(parent)
	grandParentChildren := setup.Tool.FlatChildren(grandParent)
	parentLen := len(parentChildren.GetIndices())
	grandParentLen := len(grandParentChildren.GetIndices())
	if parentLen+parentCount != grandParentLen {
		b.Errorf(
			"flat children count of parent and grand parent doesn't match. expected %v and got %v",
			parentLen+parentCount,
			grandParentLen,
		)
	}
}
