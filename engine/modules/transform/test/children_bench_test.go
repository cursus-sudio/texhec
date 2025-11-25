package test

import (
	"engine/modules/transform"
	"testing"
)

func BenchmarkChildren_1(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()

	parentTransform := setup.Transaction.GetObject(parent)
	childTransform := setup.Transaction.GetObject(child)

	childTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		parentTransform.Children()
	}
}

func BenchmarkChildren_10(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentTransform := setup.Transaction.GetObject(parent)

	for i := 0; i < 10; i++ {
		child := setup.World.NewEntity()
		childTransform := setup.Transaction.GetObject(child)
		childTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))
	}

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		parentTransform.Children()
	}
}

func BenchmarkChildren_100(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentTransform := setup.Transaction.GetObject(parent)

	for i := 0; i < 100; i++ {
		child := setup.World.NewEntity()
		childTransform := setup.Transaction.GetObject(child)
		childTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))
	}

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		parentTransform.Children()
	}
}

func BenchmarkFlatChildren_1_1(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	parentTransform := setup.Transaction.GetObject(parent)
	childTransform := setup.Transaction.GetObject(child)
	grandChildTransform := setup.Transaction.GetObject(grandChild)

	childTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))
	grandChildTransform.Parent().Set(transform.NewParent(child, transform.RelativePos))

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		parentTransform.FlatChildren()
	}
}

func BenchmarkFlatChildren_10_10(b *testing.B) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentTransform := setup.Transaction.GetObject(parent)

	for i := 0; i < 10; i++ {
		child := setup.World.NewEntity()
		childTransform := setup.Transaction.GetObject(child)
		childTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))

		for j := 0; j < 10; j++ {
			grandChild := setup.World.NewEntity()
			grandChildTransform := setup.Transaction.GetObject(grandChild)
			grandChildTransform.Parent().Set(transform.NewParent(child, transform.RelativePos))
		}
	}

	if err := setup.Transaction.Flush(); err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		parentTransform.FlatChildren()
	}
}
