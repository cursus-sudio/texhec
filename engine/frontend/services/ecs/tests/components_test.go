package ecs_test

import (
	"frontend/services/ecs"
	"testing"
)

type Component struct {
	Counter int
}

func TestComponents(t *testing.T) {
	world := ecs.NewWorld()

	component := Component{Counter: 7}
	secondComponent := Component{Counter: 8}

	if _, err := ecs.GetComponent[Component](world, ecs.EntityID(0)); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when retrieving component from not existing entity do not got ErrEntityDoNotExists error")
	}

	if err := world.SaveComponent(ecs.EntityID(0), component); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when trying to save component on not existing entity do not got ErrEntityDoNotExists error")
	}

	entityId := world.NewEntity()
	if err := world.SaveComponent(entityId, component); err != nil {
		t.Errorf("when trying to save component on existing entity got unexpected error")
	}

	if retrievedComponent, err := ecs.GetComponent[Component](world, entityId); err != nil {
		t.Errorf("unexpected error when retrieving component")
	} else if retrievedComponent != component {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	if err := world.SaveComponent(ecs.EntityID(0), secondComponent); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when trying to save existing component on not existing entity do not got ErrEntityDoNotExists error")
	}

	if err := world.SaveComponent(entityId, secondComponent); err != nil {
		t.Errorf("when saving component got unexpected error")
	}

	if retrievedComponent, err := ecs.GetComponent[Component](world, entityId); err != nil {
		t.Errorf("unexpected error when retrieving component")
	} else if retrievedComponent != secondComponent {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	world.RemoveComponent(entityId, ecs.GetComponentType(Component{}))

	if _, err := ecs.GetComponent[Component](world, entityId); err != ecs.ErrComponentDoNotExists {
		t.Errorf("retrieving removed component didn't return ecs.ErrComponentDoNotExists")
	}

	world.SaveComponent(entityId, component)
	world.RemoveEntity(entityId)

	if _, err := ecs.GetComponent[Component](world, entityId); err != ecs.ErrEntityDoNotExists {
		t.Errorf("retrieving component from removed entity didn't return ecs.ErrEntityDoNotExists")
	}
}
