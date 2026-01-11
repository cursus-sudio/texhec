package test

import (
	"engine/modules/transform"
	"testing"
)

func TestPivot(t *testing.T) {
	setup := NewSetup()
	entity := setup.NewEntity()

	setup.transform.Size().Set(entity, transform.NewSize(10, 10, 10))
	setup.expectAbsolutePos(t, entity, transform.NewPos(0, 0, 0))

	setup.transform.PivotPoint().Set(entity, transform.NewPivotPoint(0, 0, 0))
	setup.expectAbsolutePos(t, entity, transform.NewPos(5, 5, 5))

	setup.transform.PivotPoint().Set(entity, transform.NewPivotPoint(1, 1, 1))
	setup.expectAbsolutePos(t, entity, transform.NewPos(-5, -5, -5))
}
