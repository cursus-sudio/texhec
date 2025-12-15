package test

import (
	"testing"
)

func TestChildren(t *testing.T) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	setup.Tool.SetParent(child, parent)
	setup.Tool.SetParent(grandChild, child)

	if children := setup.Tool.Children(parent); !children.Get(child) || len(children.GetIndices()) != 1 {
		t.Errorf("expected parent to have one child %v", children.GetIndices())
		return
	}

	if children := setup.Tool.FlatChildren(parent); !children.Get(child) || !children.Get(grandChild) || len(children.GetIndices()) != 2 {
		t.Errorf("expected parent to have two flat children %v", children.GetIndices())
		return
	}
}
