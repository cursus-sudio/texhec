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

	var expected []ecs.EntityID
	child1 := setup.World.NewEntity()
	child2 := setup.World.NewEntity()

	expected = []ecs.EntityID{child1, child2}
	setup.Tool.SetChildren(parent, expected...)
	if children := setup.Tool.Children(parent).GetIndices(); !slices.Equal(children, expected) {
		t.Errorf("setChildren doesn't work expects %v and has %v", expected, children)
		return
	}

	expected = []ecs.EntityID{child2, child1}
	setup.Tool.SetChildren(parent, expected...)
	if children := setup.Tool.Children(parent).GetIndices(); !slices.Equal(children, expected) {
		t.Errorf("setChildren doesn't work expects %v and has %v", expected, children)
		return
	}

	expected = []ecs.EntityID{child1, child2}
	setup.Tool.SetChildren(parent, expected...)
	if children := setup.Tool.Children(parent).GetIndices(); !slices.Equal(children, expected) {
		t.Errorf("setChildren doesn't work expects %v and has %v", expected, children)
		return
	}
}
