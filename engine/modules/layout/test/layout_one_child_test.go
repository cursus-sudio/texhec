package test

import (
	"engine/modules/layout"
	"engine/modules/transform"
	"testing"
)

func TestLayoutForOneChild(t *testing.T) {
	setup := NewSetup(t)

	parent := setup.World.NewEntity()

	setup.Layout.Order().Set(parent, layout.NewOrder(layout.OrderVectical))
	setup.Layout.Align().Set(parent, layout.NewAlign(.5, .5))
	setup.Layout.Gap().Set(parent, layout.NewGap(10))

	btn := setup.World.NewEntity()

	setup.Hierarchy.SetParent(btn, parent)

	setup.Transform.Parent().Set(btn, transform.NewParent(transform.RelativePos))
	setup.Transform.Size().Set(btn, transform.NewSize(10, 10, 10))
	setup.Transform.Size().Set(parent, transform.NewSize(10, 10, 10))
	setup.Expect(btn, 0, 0)

	setup.Layout.Align().Set(parent, layout.NewAlign(1, 1))
	setup.Expect(btn, 0, 0)

	setup.Layout.Align().Set(parent, layout.NewAlign(0, 0))
	setup.Expect(btn, 0, 0)

	setup.Transform.Size().Set(btn, transform.NewSize(0, 0, 0))
	setup.Expect(btn, -5, 5)

	setup.Layout.Align().Set(parent, layout.NewAlign(1, 1))
	setup.Expect(btn, 5, -5)
}
