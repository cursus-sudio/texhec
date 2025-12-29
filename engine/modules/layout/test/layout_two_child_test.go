package test

import (
	"engine/modules/layout"
	"engine/modules/transform"
	"testing"
)

func TestLayoutForTwoChild(t *testing.T) {
	setup := NewSetup(t)

	parent := setup.NewEntity()

	setup.Transform().Size().Set(parent, transform.NewSize(0, 0, 0))

	setup.Layout().Order().Set(parent, layout.NewOrder(layout.OrderVectical))
	setup.Layout().Align().Set(parent, layout.NewAlign(.5, .5))
	setup.Layout().Gap().Set(parent, layout.NewGap(10))

	//

	btn1 := setup.NewEntity()
	setup.Transform().Size().Set(btn1, transform.NewSize(0, 0, 0))
	setup.Hierarchy().SetParent(btn1, parent)
	setup.Transform().Parent().Set(btn1, transform.NewParent(transform.RelativePos))

	btn2 := setup.NewEntity()
	setup.Transform().Size().Set(btn2, transform.NewSize(0, 0, 0))
	setup.Hierarchy().SetParent(btn2, parent)
	setup.Transform().Parent().Set(btn2, transform.NewParent(transform.RelativePos))

	setup.Expect(btn1, 0, 5)
	setup.Expect(btn2, 0, -5)
	setup.RemoveEntity(parent)
}
