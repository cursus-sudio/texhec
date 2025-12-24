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

	entity := setup.NewEntity()

	setup.Hierarchy().SetParent(entity, parent)
	setup.Transform().Parent().Set(entity, transform.NewParent(transform.RelativePos))
	setup.expectAbsolutePos(entity, transform.NewPos(10, 10, 10))

	setup.Transform().ParentPivotPoint().Set(entity, transform.NewParentPivotPoint(0, 0, 0))
	setup.expectAbsolutePos(entity, transform.NewPos(5, 5, 5))

	setup.Transform().ParentPivotPoint().Set(entity, transform.NewParentPivotPoint(1, 1, 1))
	setup.expectAbsolutePos(entity, transform.NewPos(15, 15, 15))
}
