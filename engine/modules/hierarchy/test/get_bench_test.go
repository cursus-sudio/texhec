package test

import (
	"engine/modules/hierarchy"
	"testing"
)

func BenchmarkChildren_1(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()

	parentObj := setup.Transaction.GetObject(parent)
	childObj := setup.Transaction.GetObject(child)

	childObj.Parent().Set(hierarchy.NewParent(parent))

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parentObj.Children()
	}
}

func BenchmarkChildren_10(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentObj := setup.Transaction.GetObject(parent)

	for i := 0; i < 10; i++ {
		child := setup.World.NewEntity()
		object := setup.Transaction.GetObject(child)
		object.Parent().Set(hierarchy.NewParent(parent))
	}

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parentObj.Children()
	}
}

func BenchmarkChildren_100(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentObj := setup.Transaction.GetObject(parent)

	for i := 0; i < 100; i++ {
		child := setup.World.NewEntity()
		object := setup.Transaction.GetObject(child)
		object.Parent().Set(hierarchy.NewParent(parent))
	}

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parentObj.Children()
	}
}

func BenchmarkFlatChildren_1_1(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	parentObj := setup.Transaction.GetObject(parent)
	childObj := setup.Transaction.GetObject(child)
	grandChildObj := setup.Transaction.GetObject(grandChild)

	childObj.Parent().Set(hierarchy.NewParent(parent))
	grandChildObj.Parent().Set(hierarchy.NewParent(child))

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parentObj.FlatChildren()
	}
}

func BenchmarkFlatChildren_10_10(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentObj := setup.Transaction.GetObject(parent)

	for i := 0; i < 10; i++ {
		child := setup.World.NewEntity()
		childObj := setup.Transaction.GetObject(child)
		childObj.Parent().Set(hierarchy.NewParent(parent))

		for j := 0; j < 10; j++ {
			grandChild := setup.World.NewEntity()
			grandChildObj := setup.Transaction.GetObject(grandChild)
			grandChildObj.Parent().Set(hierarchy.NewParent(child))
		}
	}

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		parentObj.FlatChildren()
	}
}
