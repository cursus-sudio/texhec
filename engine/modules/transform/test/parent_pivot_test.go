package test

import (
	"engine/modules/transform"
	"testing"
)

func TestParentPivot(t *testing.T) {
	setup := NewSetup()
	parent := setup.NewEntity()
	setup.transform.Pos().Set(parent, transform.NewPos(10, 10, 10))
	setup.transform.Size().Set(parent, transform.NewSize(10, 10, 10))

	child := setup.NewEntity()

	setup.hierarchy.SetParent(child, parent)
	setup.transform.Parent().Set(child, transform.NewParent(transform.RelativePos))
	setup.expectAbsolutePos(t, child, transform.NewPos(10, 10, 10))

	setup.transform.ParentPivotPoint().Set(child, transform.NewParentPivotPoint(0, 0, 0))
	setup.expectAbsolutePos(t, child, transform.NewPos(5, 5, 5))

	setup.transform.ParentPivotPoint().Set(child, transform.NewParentPivotPoint(1, 1, 1))
	setup.expectAbsolutePos(t, child, transform.NewPos(15, 15, 15))

	setup.transform.ParentPivotPoint().Set(child, transform.NewParentPivotPoint(0, 0, 0))
	setup.transform.PivotPoint().Set(child, transform.NewPivotPoint(0, 0, 0))
	setup.transform.Size().Set(child, transform.NewSize(0, 0, 0))
	setup.expectAbsolutePos(t, child, transform.NewPos(5, 5, 5))
	setup.expectAbsolutePos(t, child, transform.NewPos(5, 5, 5))
}
