package ecs

import (
	// "reflect"
	"math"
)

// interface

type ComponentsArray[Component any] interface {
	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) error // upsert

	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetComponent(entity EntityID) (Component, error)
	RemoveComponent(EntityID)

	OnAdd(func([]EntityID))
	OnChange(func([]EntityID))
	OnRemove(func([]EntityID))
}

// impl

const (
	noEntity    uint32 = math.MaxUint32 - 0
	noComponent uint32 = math.MaxUint32 - 1
)

type componentsArray[Component any] struct {
	// entity here is an index
	entitiesComponents []uint32 // here some indices have special meaning (read constants above)
	components         []Component
	componentsEntities []EntityID

	onRemove []func([]EntityID)
	onChange []func([]EntityID)
	onAdd    []func([]EntityID)
}

func NewComponentsArray[Component any]() *componentsArray[Component] {
	return &componentsArray[Component]{
		entitiesComponents: make([]uint32, 0),
		components:         make([]Component, 0),
		componentsEntities: make([]EntityID, 0),

		onRemove: make([]func([]EntityID), 0),
		onChange: make([]func([]EntityID), 0),
		onAdd:    make([]func([]EntityID), 0),
	}
}

func (c *componentsArray[Component]) SaveComponent(entity EntityID, component Component) error {
	listener := func() {}
	defer listener()

	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return ErrEntityDoNotExists
	}
	componentIndex := c.entitiesComponents[entityIndex]
	if componentIndex == noEntity {
		return ErrEntityDoNotExists
	}
	// component := componentsArrayComponent[Component]{entity, rawComponent}
	if componentIndex == noComponent {
		c.entitiesComponents[entityIndex] = uint32(len(c.components))
		c.components = append(c.components, component)
		c.componentsEntities = append(c.componentsEntities, entity)
		listener = func() {
			entities := []EntityID{entity}
			for _, listener := range c.onAdd {
				listener(entities)
			}
		}
		return nil
	}
	c.components[componentIndex] = component
	c.componentsEntities[componentIndex] = entity
	listener = func() {
		entities := []EntityID{entity}
		for _, listener := range c.onChange {
			listener(entities)
		}
	}
	return nil
}

func (c *componentsArray[Component]) GetComponent(entity EntityID) (Component, error) {
	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		var zero Component
		return zero, ErrEntityDoNotExists
	}
	componentIndex := c.entitiesComponents[entityIndex]
	var zero Component
	switch componentIndex {
	case noEntity:
		return zero, ErrComponentDoNotExists
	case noComponent:
		return zero, ErrEntityDoNotExists
	default:
		component := c.components[componentIndex]
		return component, nil
	}
}

func (c *componentsArray[Component]) RemoveComponent(entity EntityID) {
	listener := func() {}
	defer listener()

	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return
	}
	componentIndex := c.entitiesComponents[entityIndex]
	if componentIndex == noEntity || componentIndex == noComponent {
		return
	}

	c.entitiesComponents[entityIndex] = noComponent

	// move component
	movedComponent := c.components[len(c.components)-1]
	c.components[componentIndex] = movedComponent
	c.components = c.components[:len(c.components)-1]

	movedComponentEntity := c.componentsEntities[len(c.componentsEntities)-1]
	c.componentsEntities[componentIndex] = movedComponentEntity
	c.componentsEntities = c.componentsEntities[:len(c.componentsEntities)-1]

	movedEntityIndex := movedComponentEntity.Index()
	c.entitiesComponents[movedEntityIndex] = componentIndex

	listener = func() {
		entities := []EntityID{entity}
		for _, listener := range c.onRemove {
			listener(entities)
		}
	}
}

func (c *componentsArray[Component]) AddEntity(entity EntityID) {

	entityIndex := entity.Index()
	for entityIndex >= len(c.entitiesComponents) {
		c.entitiesComponents = append(c.entitiesComponents, noEntity)
	}
	c.entitiesComponents[entityIndex] = noComponent
}

func (c *componentsArray[Component]) RemoveEntity(entity EntityID) {
	listener := func() {}
	defer listener()

	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return
	}
	componentIndex := c.entitiesComponents[entityIndex]
	if componentIndex != noComponent && componentIndex != noEntity {
		listener = func() {
			entities := []EntityID{entity}
			for _, listener := range c.onRemove {
				listener(entities)
			}
		}
	}
	c.entitiesComponents[entityIndex] = noEntity
}

func (c *componentsArray[Component]) OnAdd(listener func([]EntityID)) {
	c.onAdd = append(c.onAdd, listener)
}

func (c *componentsArray[Component]) OnChange(listener func([]EntityID)) {
	c.onChange = append(c.onChange, listener)
}

func (c *componentsArray[Component]) OnRemove(listener func([]EntityID)) {
	c.onRemove = append(c.onRemove, listener)
}
