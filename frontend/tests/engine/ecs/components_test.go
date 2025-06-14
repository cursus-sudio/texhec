package ecs_test

import (
	"frontend/src/engine/ecs"
	"testing"

	"github.com/ogiusek/ioc"
)

type Component struct {
	Counter int
}

func TestComponents(t *testing.T) {
	c := ioc.NewContainer()
	ecs.Package().Register(c)
	world := ioc.Get[ecs.WorldFactory](c)()

	component := Component{Counter: 7}
	secondComponent := Component{Counter: 8}

	componentType := ecs.GetComponentType(component)
	if err := world.GetComponent(ecs.EntityId{}, &componentType); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when retrieving component from not existing entity do not got ErrEntityDoNotExists error")
	}

	if err := world.SaveComponent(ecs.EntityId{}, component); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when trying to save component on not existing entity do not got ErrEntityDoNotExists error")
	}

	entityId := world.NewEntity()
	if err := world.SaveComponent(entityId, component); err != nil {
		t.Errorf("when trying to save component on existing entity got unexpected error")
	}

	retrievedComponent := Component{}
	if err := world.GetComponent(entityId, &retrievedComponent); err != nil {
		t.Errorf("unexpected error when retrieving component")
	}

	if retrievedComponent != component {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	if err := world.SaveComponent(ecs.EntityId{}, secondComponent); err != ecs.ErrEntityDoNotExists {
		t.Errorf("when trying to save existing component on not existing entity do not got ErrEntityDoNotExists error")
	}

	if err := world.SaveComponent(entityId, secondComponent); err != nil {
		t.Errorf("when saving component got unexpected error")
	}

	if err := world.GetComponent(entityId, &retrievedComponent); err != nil {
		t.Errorf("unexpected error when retrieving component")
	}

	if retrievedComponent != secondComponent {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	world.RemoveComponent(entityId, componentType)

	if err := world.GetComponent(entityId, &retrievedComponent); err != ecs.ErrComponentDoNotExists {
		t.Errorf("retrieving removed component didn't return ecs.ErrComponentDoNotExists")
	}

	world.SaveComponent(entityId, component)
	world.RemoveEntity(entityId)

	if err := world.GetComponent(entityId, &retrievedComponent); err != ecs.ErrEntityDoNotExists {
		t.Errorf("retrieving component from removed entity didn't return ecs.ErrEntityDoNotExists")
	}
}
