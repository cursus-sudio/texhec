package test

import (
	"engine/modules/transform"
	"testing"
)

func TestParentPivot(t *testing.T) {
	setup := NewSetup(t)
	parent := setup.NewEntity()
	setup.Transform().Pos().Set(parent, transform.NewPos(10, 10, 10))
	setup.Transform().Size().Set(parent, transform.NewSize(10, 10, 10))

	child := setup.NewEntity()

	setup.Hierarchy().SetParent(child, parent)
	setup.Transform().Parent().Set(child, transform.NewParent(transform.RelativePos))
	setup.expectAbsolutePos(child, transform.NewPos(10, 10, 10))

	setup.Transform().ParentPivotPoint().Set(child, transform.NewParentPivotPoint(0, 0, 0))
	setup.expectAbsolutePos(child, transform.NewPos(5, 5, 5))

	setup.Transform().ParentPivotPoint().Set(child, transform.NewParentPivotPoint(1, 1, 1))
	setup.expectAbsolutePos(child, transform.NewPos(15, 15, 15))

	setup.Transform().ParentPivotPoint().Set(child, transform.NewParentPivotPoint(0, 0, 0))
	setup.Transform().PivotPoint().Set(child, transform.NewPivotPoint(0, 0, 0))
	setup.Transform().Size().Set(child, transform.NewSize(0, 0, 0))
	setup.expectAbsolutePos(child, transform.NewPos(5, 5, 5))
	setup.expectAbsolutePos(child, transform.NewPos(5, 5, 5))
}
