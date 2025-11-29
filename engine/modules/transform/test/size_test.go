package test

import (
	"engine/modules/transform"
	"testing"
)

func TestAbsoluteSize(t *testing.T) {
	setup := NewSetup()
	entity := setup.World.NewEntity()

	expectSize := func(expectedSize transform.SizeComponent) {
		entityTransform := setup.Transaction.GetObject(entity)
		size, err := entityTransform.AbsoluteSize().Get()
		if err != nil {
			t.Error(err)
			return
		}
		if size != expectedSize {
			t.Errorf("expected size %v but has %v", expectedSize, size)
		}
	}

	{
		entityTransform := setup.Transaction.GetObject(entity)
		entityTransform.Size().Set(transform.NewSize(10, 10, 10))
		if err := setup.Transaction.Flush(); err != nil {
			t.Error(err)
			return
		}
		expectSize(transform.NewSize(10, 10, 10))
	}

	{
		entityTransform := setup.Transaction.GetObject(entity)
		entityTransform.Size().Set(transform.NewSize(15, 15, 15))
		if err := setup.Transaction.Flush(); err != nil {
			t.Error(err)
			return
		}
		expectSize(transform.NewSize(15, 15, 15))
	}
}
