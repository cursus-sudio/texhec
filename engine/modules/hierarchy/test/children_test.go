package test

import (
	"engine/services/ecs"
	"slices"
	"testing"
)

func TestChildren(t *testing.T) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	setup.Tool.SetParent(child, parent)
	setup.Tool.SetParent(grandChild, child)

	if parents := setup.Tool.GetOrderedParents(child); len(parents) != 1 || parents[0] != parent {
		t.Errorf("expected [%v] parents of a child but has %v", parent, parents)
		return
	}

	if parents := setup.Tool.GetOrderedParents(grandChild); len(parents) != 2 || parents[0] != child || parents[1] != parent {
		t.Errorf("expected [%v %v] parents of a grand child but has %v", child, parent, parents)
		return
	}

	if children := setup.Tool.Children(parent); !children.Get(child) || len(children.GetIndices()) != 1 {
		t.Errorf("expected parent to have one child %v", children.GetIndices())
		return
	}

	if children := setup.Tool.FlatChildren(parent); !children.Get(child) || !children.Get(grandChild) || len(children.GetIndices()) != 2 {
		t.Errorf("expected parent to have two flat children %v", children.GetIndices())
		return
	}

	setup.World.RemoveEntity(parent)
	setup.Tool.Children(parent)
	if exists := setup.World.EntityExists(child); exists {
		t.Errorf("parent children should be removed")
		return
	}
	if exists := setup.World.EntityExists(grandChild); exists {
		t.Errorf("parent children should be removed")
		return
	}
	if children := setup.Tool.Children(parent); len(children.GetIndices()) != 0 {
		t.Errorf("removed entity still has children")
		return
	}
	if children := setup.Tool.FlatChildren(parent); len(children.GetIndices()) != 0 {
		t.Errorf("removed entity still has children")
		return
	}
	if children := setup.Tool.Children(parent); len(children.GetIndices()) != 0 {
		t.Errorf("removed entity still has children")
		return
	}
	if children := setup.Tool.FlatChildren(parent); len(children.GetIndices()) != 0 {
		t.Errorf("removed entity still has children")
		return
	}
}

func TestSetChildren(t *testing.T) {
	setup := NewSetup()

	parent := setup.World.NewEntity()

	setAndExpect := func(expected ...ecs.EntityID) {
		t.Helper()
		setup.Tool.SetChildren(parent, expected...)
		if children := setup.Tool.Children(parent).GetIndices(); !slices.Equal(children, expected) {
			t.Errorf("setChildren doesn't work expects %v and has %v", expected, children)
		}
	}

	c1 := setup.World.NewEntity()
	c2 := setup.World.NewEntity()
	c3 := setup.World.NewEntity()
	c4 := setup.World.NewEntity()

	// this order tests does removal works for more entities
	setAndExpect(c1, c2, c3, c4)
	setAndExpect(c2, c1, c4, c3)
	setAndExpect(c1, c2, c3, c4)

	setAndExpect(c1, c2)
	setAndExpect(c2, c1)
	setAndExpect(c1, c2)

	setAndExpect(10010, 10011, 10012, 10013, 10014) // this is real example
	setAndExpect(10010, 10011, 10012, 10013, 10014)
	setAndExpect(10010, 10013, 10012, 10011, 10014)

}
