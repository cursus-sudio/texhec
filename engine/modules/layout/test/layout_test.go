package test

import (
	"engine/modules/layout"
	"engine/modules/transform"
	"engine/services/ecs"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestLayout(t *testing.T) {
	setup := NewSetup()

	if true {
		// this test isn't finished
		return
	}

	parent := setup.NewEntity()
	// setup.Transform().Size().Set(parent, transform.NewSize(500, 200, 1))
	// setup.Hierarchy().SetParent(buttonArea, cameraEntity)
	// setup.Transform().Parent().Set(parent, transform.NewParent(transform.RelativePos))

	setup.Layout().Order().Set(parent, layout.NewOrder(layout.OrderVectical))
	setup.Layout().Align().Set(parent, layout.NewAlign(.5, .5))
	setup.Layout().Gap().Set(parent, layout.NewGap(10))

	//

	btns := []ecs.EntityID{}
	for i := 0; i < 2; i++ {
		btns = append(btns, setup.NewEntity())
	}

	for _, btnEntity := range btns {
		setup.Hierarchy().SetParent(btnEntity, parent)

		setup.Transform().AspectRatio().Set(btnEntity, transform.NewAspectRatio(1, 1, 0, transform.PrimaryAxisX))
		// setup.Transform().Parent().Set(btnEntity, transform.NewParent(transform.RelativePos|transform.RelativeSizeX))
		setup.Transform().Parent().Set(btnEntity, transform.NewParent(transform.RelativePos))
		// setup.Transform().MaxSize().Set(btnEntity, transform.NewMaxSize(0, 50, 0))
		setup.Transform().Size().Set(btnEntity, transform.NewSize(1, 50, 1))
	}

	if pos, _ := setup.Transform().AbsolutePos().Get(btns[0]); pos.Pos[2] != 30 {
		t.Errorf("expected %v and got %v", mgl32.Vec3{0, 30, 0}, pos)
	}
	if pos, _ := setup.Transform().AbsolutePos().Get(btns[1]); pos.Pos[2] != -30 {
		t.Errorf("expected %v and got %v", mgl32.Vec3{0, -30, 0}, pos)
	}
}
