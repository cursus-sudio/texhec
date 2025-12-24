package ecs_test

import (
	"engine/services/ecs"
	"testing"
)

type Component struct {
	Counter int
}

var component = Component{Counter: 7}
var secondComponent = Component{Counter: 8}

func TestComponents(t *testing.T) {
	world := ecs.NewWorld()

	if _, ok := ecs.GetComponent[Component](world, ecs.EntityID(0)); ok {
		t.Errorf("retrieved not existing component")
	}

	entityId := world.NewEntity()
	ecs.SaveComponent(world, entityId, component)

	if retrievedComponent, ok := ecs.GetComponent[Component](world, entityId); !ok {
		t.Errorf("expected component")
	} else if retrievedComponent != component {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	ecs.SaveComponent(world, entityId, secondComponent)

	if retrievedComponent, ok := ecs.GetComponent[Component](world, entityId); !ok {
		t.Errorf("expected component")
	} else if retrievedComponent != secondComponent {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	ecs.RemoveComponent[Component](world, entityId)

	if _, ok := ecs.GetComponent[Component](world, entityId); ok {
		t.Errorf("retrieved removed component")
	}

	ecs.SaveComponent(world, entityId, component)
	world.RemoveEntity(entityId)

	if _, ok := ecs.GetComponent[Component](world, entityId); ok {
		t.Errorf("retrieved removed component")
	}
}

func TestComponentsArrays(t *testing.T) {
	world := ecs.NewWorld()
	componentArray := ecs.GetComponentsArray[Component](world)

	entityId := world.NewEntity()
	componentArray.Set(entityId, component)

	if retrievedComponent, ok := componentArray.Get(entityId); !ok {
		t.Errorf("expected component")
	} else if retrievedComponent != component {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	componentArray.Set(entityId, secondComponent)

	if retrievedComponent, ok := componentArray.Get(entityId); !ok {
		t.Errorf("expected component")
	} else if retrievedComponent != secondComponent {
		t.Errorf("retrieved component isn't equal to saved component")
	}

	componentArray.Remove(entityId)

	if _, ok := componentArray.Get(entityId); ok {
		t.Errorf("retrieved removed component")
	}

	componentArray.Set(entityId, component)
	world.RemoveEntity(entityId)

	if _, ok := componentArray.Get(entityId); ok {
		t.Errorf("retrieved removed component")
	}
}

func TestComponentsQuery(t *testing.T) {
	type Component2 struct{}
	world := ecs.NewWorld()

	component := ecs.GetComponentsArray[Component](world)
	component2 := ecs.GetComponentsArray[Component2](world)

	set := ecs.NewDirtySet()
	component.AddDirtySet(set)
	component2.AddDirtySet(set)

	if dirty := set.Get(); len(dirty) != 0 {
		t.Errorf("no dirty flags were expected")
		return
	}

	entity := world.NewEntity()

	component.Set(entity, Component{})
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected entity to be dirty")
		return
	}

	component.Remove(entity)
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected entity to be dirty")
		return
	}

	component2.Set(entity, Component2{})
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected entity to be dirty")
		return
	}

	component2.Remove(entity)
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected entity to be dirty")
		return
	}
}
