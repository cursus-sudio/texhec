package ecs_test

import (
	"shared/services/ecs"
	"testing"
)

type Component struct {
	Counter int
}

var component = Component{Counter: 7}
var secondComponent = Component{Counter: 8}

func TestComponents(t *testing.T) {
	world := ecs.NewWorld()

	if _, err := ecs.GetComponent[Component](world.Components(), ecs.EntityID(0)); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when retrieving component from not existing entity do not got ErrEntityDoNotExists error")
	}

	entityId := world.NewEntity()
	if err := ecs.SaveComponent(world.Components(), entityId, component); err != nil {
		t.Errorf("when trying to save component on existing entity got unexpected error")
	}

	if retrievedComponent, err := ecs.GetComponent[Component](world.Components(), entityId); err != nil {
		t.Errorf("unexpected error when retrieving component")
	} else if retrievedComponent != component {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	if err := ecs.SaveComponent(world.Components(), ecs.EntityID(0), secondComponent); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when trying to save existing component on not existing entity do not got ErrEntityDoNotExists error")
	}

	if err := ecs.SaveComponent(world.Components(), entityId, secondComponent); err != nil {
		t.Errorf("when saving component got unexpected error")
	}

	if retrievedComponent, err := ecs.GetComponent[Component](world.Components(), entityId); err != nil {
		t.Errorf("unexpected error when retrieving component")
	} else if retrievedComponent != secondComponent {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	ecs.RemoveComponent[Component](world.Components(), entityId)

	if _, err := ecs.GetComponent[Component](world.Components(), entityId); err != ecs.ErrComponentDoNotExists {
		t.Errorf("retrieving removed component didn't return ecs.ErrComponentDoNotExists but %v\n", err)
	}

	ecs.SaveComponent(world.Components(), entityId, component)
	world.RemoveEntity(entityId)

	if _, err := ecs.GetComponent[Component](world.Components(), entityId); err != ecs.ErrEntityDoNotExists {
		t.Errorf("retrieving component from removed entity didn't return ecs.ErrEntityDoNotExists but %v\n", err)
	}
}

func TestComponentsArrays(t *testing.T) {
	world := ecs.NewWorld()
	componentArray := ecs.GetComponentsArray[Component](world.Components())

	if _, err := componentArray.GetComponent(ecs.EntityID(0)); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when retrieving component from not existing entity do not got ErrEntityDoNotExists error")
	}

	entityId := world.NewEntity()
	if err := componentArray.SaveComponent(entityId, component); err != nil {
		t.Errorf("when trying to save component on existing entity got unexpected error")
	}

	if retrievedComponent, err := componentArray.GetComponent(entityId); err != nil {
		t.Errorf("unexpected error when retrieving component")
	} else if retrievedComponent != component {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	if err := componentArray.SaveComponent(ecs.EntityID(0), secondComponent); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when trying to save existing component on not existing entity do not got ErrEntityDoNotExists error")
	}

	if err := componentArray.SaveComponent(entityId, secondComponent); err != nil {
		t.Errorf("when saving component got unexpected error")
	}

	if retrievedComponent, err := componentArray.GetComponent(entityId); err != nil {
		t.Errorf("unexpected error when retrieving component")
	} else if retrievedComponent != secondComponent {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	componentArray.RemoveComponent(entityId)

	if _, err := componentArray.GetComponent(entityId); err != ecs.ErrComponentDoNotExists {
		t.Errorf("retrieving removed component didn't return ecs.ErrComponentDoNotExists but %v\n", err)
	}

	componentArray.SaveComponent(entityId, component)
	world.RemoveEntity(entityId)

	if _, err := componentArray.GetComponent(entityId); err != ecs.ErrEntityDoNotExists {
		t.Errorf("retrieving component from removed entity didn't return ecs.ErrEntityDoNotExists but %v\n", err)
	}
}

func TestComponentsQuery(t *testing.T) {
	type Component2 struct{}
	type ForbiddenComponent struct{}
	type TrackedComponent struct{}
	world := ecs.NewWorld()

	adds := 0
	expectedAdds := 0
	changes := 0
	expectedChanges := 0
	removes := 0
	expectedRemoves := 0

	query := world.Query().
		Require(ecs.GetComponentType(Component{})).
		Require(ecs.GetComponentType(Component2{})).
		Forbid(ecs.GetComponentType(ForbiddenComponent{})).
		Track(ecs.GetComponentType(TrackedComponent{})).
		Build()

	query.OnAdd(func(ei []ecs.EntityID) { adds += 1 })
	query.OnChange(func(ei []ecs.EntityID) { changes += 1 })
	query.OnRemove(func(ei []ecs.EntityID) { removes += 1 })

	entity := world.NewEntity()

	component := ecs.GetComponentsArray[Component](world.Components())
	component2 := ecs.GetComponentsArray[Component2](world.Components())
	forbiddenComponent := ecs.GetComponentsArray[ForbiddenComponent](world.Components())
	trackedComponent := ecs.GetComponentsArray[TrackedComponent](world.Components())

	expectNothing := func() bool {
		if adds != expectedAdds {
			t.Errorf("unexpected call on query onAdd")
			return false
		}
		if changes != expectedChanges {
			t.Errorf("unexpected call on query onChange")
			return false
		}
		if removes != expectedRemoves {
			t.Errorf("unexpected call on query onRemove")
			return false
		}
		return true
	}

	expectAdd := func() bool {
		expectedAdds += 1
		if adds != expectedAdds {
			t.Errorf("expected call on query onAdd")
			return false
		}
		if changes != expectedChanges {
			t.Errorf("unexpected call on query onChange expected call onAdd")
			return false
		}
		if removes != expectedRemoves {
			t.Errorf("unexpected call on query onRemove expected call onAdd")
			return false
		}
		return true
	}

	expectChange := func() bool {
		expectedChanges += 1
		if adds != expectedAdds {
			t.Errorf("unexpected call on query onAdd expected onChange")
			return false
		}
		if changes != expectedChanges {
			t.Errorf("expected call on query onChange")
			return false
		}
		if removes != expectedRemoves {
			t.Errorf("unexpected call on query onRemove expected call onChange")
			return false
		}
		return true
	}

	expectRemove := func() bool {
		expectedRemoves += 1
		if adds != expectedAdds {
			t.Errorf("unexpected call on query onAdd expected onRemove")
			return false
		}
		if changes != expectedChanges {
			t.Errorf("unexpected call on query onChange expected onRemove")
			return false
		}
		if removes != expectedRemoves {
			t.Errorf("expected call on query onRemove")
			return false
		}
		return true
	}

	component.SaveComponent(entity, Component{})
	if ok := expectNothing(); !ok {
		return
	}

	component2.SaveComponent(entity, Component2{})
	if ok := expectAdd(); !ok {
		return
	}

	component2.SaveComponent(entity, Component2{})
	if ok := expectChange(); !ok {
		return
	}

	component.SaveComponent(entity, Component{})
	if ok := expectChange(); !ok {
		return
	}

	trackedComponent.SaveComponent(entity, TrackedComponent{})
	if ok := expectChange(); !ok {
		return
	}

	trackedComponent.SaveComponent(entity, TrackedComponent{})
	if ok := expectChange(); !ok {
		return
	}

	trackedComponent.RemoveComponent(entity)
	if ok := expectChange(); !ok {
		return
	}

	forbiddenComponent.SaveComponent(entity, ForbiddenComponent{})
	if ok := expectRemove(); !ok {
		return
	}

	forbiddenComponent.RemoveComponent(entity)
	if ok := expectAdd(); !ok {
		return
	}

	component.RemoveComponent(entity)
	if ok := expectRemove(); !ok {
		return
	}
}
