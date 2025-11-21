package test

import (
	"frontend/modules/transform"
	"testing"
)

func TestChildren(t *testing.T) {
	setup := NewSetup()
	parent := setup.World.NewEntity()
	child := setup.World.NewEntity()
	grandChild := setup.World.NewEntity()

	parentTransform := setup.Transaction.GetEntity(parent)
	childTransform := setup.Transaction.GetEntity(child)
	grandChildTransform := setup.Transaction.GetEntity(grandChild)

	childTransform.Parent().Set(transform.NewParent(parent, transform.RelativePos))
	grandChildTransform.Parent().Set(transform.NewParent(child, transform.RelativePos))

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
