package parent_test

import (
	"testing"
)

func TestIndexing(t *testing.T) {
	setup := NewSetup()
	parent := setup.W.NewEntity()
	if children := setup.Tool().GetChildren(parent); len(children.GetIndices()) != 0 {
		t.Errorf("unexpected children")
		return
	}

	child := setup.W.NewEntity()
	component := Component{Parent: parent}
	setup.Array.SaveComponent(child, component)

	children := setup.Tool().GetChildren(parent)
	if len(children.GetIndices()) != 1 {
		t.Errorf("expected 1 entity to be added to parent")
		return
	}
	retrievedChild := children.GetIndices()[0]

	if retrievedChild != child {
		t.Errorf("invalid child added. expected %v but has %v", child, retrievedChild)
		return
	}

	{
		child := setup.W.NewEntity()
		component := Component{Parent: parent}
		setup.Array.SaveComponent(child, component)

		children := setup.Tool().GetChildren(parent)
		if len(children.GetIndices()) != 2 {
			t.Errorf("expected 1 entity to be added to parent")
			return
		}
		chilIsOkay := children.Get(child)

		if !chilIsOkay {
			t.Errorf("invalid child added. expected %v", child)
			return
		}
		setup.Array.RemoveComponent(child)
		if children := setup.Tool().GetChildren(parent); len(children.GetIndices()) != 1 {
			t.Errorf("expected children to be removed")
			return
		}
	}

	setup.Array.RemoveComponent(child)
	if children := setup.Tool().GetChildren(parent); len(children.GetIndices()) != 0 {
		t.Errorf("expected children to be removed")
		return
	}
}
