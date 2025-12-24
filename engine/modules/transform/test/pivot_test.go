package test

import (
	"engine/modules/transform"
	"testing"
)

func TestPivot(t *testing.T) {
	setup := NewSetup(t)
	entity := setup.NewEntity()

	setup.Transform().Size().Set(entity, transform.NewSize(10, 10, 10))
	setup.expectAbsolutePos(entity, transform.NewPos(0, 0, 0))

	setup.Transform().PivotPoint().Set(entity, transform.NewPivotPoint(0, 0, 0))
	setup.expectAbsolutePos(entity, transform.NewPos(5, 5, 5))

	setup.Transform().PivotPoint().Set(entity, transform.NewPivotPoint(1, 1, 1))
	setup.expectAbsolutePos(entity, transform.NewPos(-5, -5, -5))
}
