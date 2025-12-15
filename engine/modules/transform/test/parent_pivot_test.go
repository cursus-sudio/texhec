package test

import (
	"engine/modules/transform"
	"testing"
)

func TestParentPivot(t *testing.T) {
	setup := NewSetup(t)
	parent := setup.World.NewEntity()
	setup.Transform.Pos().SaveComponent(parent, transform.NewPos(10, 10, 10))
	setup.Transform.Size().SaveComponent(parent, transform.NewSize(10, 10, 10))

	entity := setup.World.NewEntity()

	setup.Hierarchy.SetParent(entity, parent)
	setup.Transform.Parent().SaveComponent(entity, transform.NewParent(transform.RelativePos))
	setup.expectAbsolutePos(entity, transform.NewPos(10, 10, 10))

	setup.Transform.ParentPivotPoint().SaveComponent(entity, transform.NewParentPivotPoint(0, 0, 0))
	setup.expectAbsolutePos(entity, transform.NewPos(5, 5, 5))

	setup.Transform.ParentPivotPoint().SaveComponent(entity, transform.NewParentPivotPoint(1, 1, 1))
	setup.expectAbsolutePos(entity, transform.NewPos(15, 15, 15))
}
