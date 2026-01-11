package test

import (
	"testing"
)

func TestIndexing(t *testing.T) {
	setup := NewSetup()
	component := Component{Index: 69}
	if _, ok := setup.Tool.Get(component.Index); ok {
		t.Errorf("expected !ok when retriving entity by not existing index")
		return
	}

	entity := setup.W.NewEntity()
	setup.Array.Set(entity, component)

	returnedEntity, ok := setup.Tool.Get(component.Index)
	if !ok {
		t.Errorf("expected ok when retriving entity by existing index")
		return
	}

	if entity != returnedEntity {
		t.Errorf("expected entities to match but they don't %v != %v\n", entity, returnedEntity)
		return
	}
}
