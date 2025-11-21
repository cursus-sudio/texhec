package test

import (
	"frontend/modules/transform"
	"testing"
)

func TestPivot(t *testing.T) {
	setup := NewSetup()
	entity := setup.World.NewEntity()
	entityTransform := setup.Transaction.GetEntity(entity)

	entityTransform.Size().Set(transform.NewSize(10, 10, 10))
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

	expectPos(transform.NewPos(0, 0, 0))

	entityTransform.PivotPoint().Set(transform.NewPivotPoint(0, 0, 0))
	if err := setup.Transaction.Flush(); err != nil {
		t.Error(err)
		return
	}
	expectPos(transform.NewPos(5, 5, 5))

	entityTransform.PivotPoint().Set(transform.NewPivotPoint(1, 1, 1))
	if err := setup.Transaction.Flush(); err != nil {
		t.Error(err)
		return
	}
	expectPos(transform.NewPos(-5, -5, -5))
}
