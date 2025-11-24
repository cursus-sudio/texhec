package test

import (
	"frontend/modules/transform"
	"testing"
)

func TestParentPivot(t *testing.T) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	parentTransform := setup.Transaction.GetObject(parent)
	parentTransform.Pos().Set(transform.NewPos(10, 10, 10))
	parentTransform.Size().Set(transform.NewSize(10, 10, 10))

	entity := setup.World.NewEntity()
	entityTransform := setup.Transaction.GetObject(entity)

	entityTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))
	if err := setup.Transaction.Flush(); err != nil {
		t.Error(err)
		return
	}

	expectPos := func(expectedPos transform.PosComponent) {
		pos, err := entityTransform.AbsolutePos().Get()
		if err != nil {
			t.Error(err)
			return
		}
		if pos != expectedPos {
			t.Errorf("expected pos %v but has %v", expectedPos, pos)
		}
	}

	expectPos(transform.NewPos(10, 10, 10))

	entityTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(0, 0, 0))
	if err := setup.Transaction.Flush(); err != nil {
		t.Error(err)
		return
	}
	expectPos(transform.NewPos(5, 5, 5))

	entityTransform.ParentPivotPoint().Set(transform.NewParentPivotPoint(1, 1, 1))
	if err := setup.Transaction.Flush(); err != nil {
		t.Error(err)
		return
	}
	expectPos(transform.NewPos(15, 15, 15))
}
