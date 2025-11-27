package test

import (
	"engine/modules/hierarchy"
	"testing"
)

func TestChildren(t *testing.T) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	parentTransform := setup.Transaction.GetObject(parent)
	childTransform := setup.Transaction.GetObject(child)
	grandChildTransform := setup.Transaction.GetObject(grandChild)

	childTransform.Parent().Set(hierarchy.NewParent(parent))
	grandChildTransform.Parent().Set(hierarchy.NewParent(child))

	if err := setup.Transaction.Flush(); err != nil {
		t.Error(err)
		return
	}

	if children := parentTransform.Children(); !children.Get(child) || len(children.GetIndices()) != 1 {
		t.Errorf("expected parent to have one child")
		return
	}

	if children := parentTransform.FlatChildren(); !children.Get(child) || !children.Get(grandChild) || len(children.GetIndices()) != 2 {
		t.Errorf("expected parent to have two flat children")
		return
	}
}
