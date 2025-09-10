package ecs

import (
	"math"
)

// interface

type ComponentsArray[Component any] interface {
	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) error // upsert

	// differs from save component by not triggering events
	// can return:
	// - ErrEntityDoNotExists
	DirtySaveComponent(EntityID, Component) error // upsert

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

	onAdd    []func([]EntityID)
	onChange []func([]EntityID)
	onRemove []func([]EntityID)
}

func NewComponentsArray[Component any]() *componentsArray[Component] {
	return &componentsArray[Component]{
		entitiesComponents: make([]uint32, 0),
		components:         make([]Component, 0),
		componentsEntities: make([]EntityID, 0),

		onAdd:    make([]func([]EntityID), 0),
		onChange: make([]func([]EntityID), 0),
		onRemove: make([]func([]EntityID), 0),
	}
}

func (c *componentsArray[Component]) SaveComponent(entity EntityID, component Component) error {
	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return ErrEntityDoNotExists
	}
	componentIndex := c.entitiesComponents[entityIndex]
	if componentIndex == noEntity {
		return ErrEntityDoNotExists
	}
	if componentIndex == noComponent {
		c.entitiesComponents[entityIndex] = uint32(len(c.components))
		c.components = append(c.components, component)
		c.componentsEntities = append(c.componentsEntities, entity)

		// listeners
		entities := []EntityID{entity}
		for _, listener := range c.onAdd {
			listener(entities)
		}
		return nil
	}
	c.components[componentIndex] = component
	c.componentsEntities[componentIndex] = entity

	// listeners
	entities := []EntityID{entity}
	for _, listener := range c.onChange {
		listener(entities)
	}

	return nil
}

func (c *componentsArray[Component]) DirtySaveComponent(entity EntityID, component Component) error {
	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return ErrEntityDoNotExists
	}
	componentIndex := c.entitiesComponents[entityIndex]
	if componentIndex == noEntity {
		return ErrEntityDoNotExists
	}
	if componentIndex == noComponent {
		c.entitiesComponents[entityIndex] = uint32(len(c.components))
		c.components = append(c.components, component)
		c.componentsEntities = append(c.componentsEntities, entity)
		return nil
	}
	c.components[componentIndex] = component
	c.componentsEntities[componentIndex] = entity
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
		return zero, ErrEntityDoNotExists
	case noComponent:
		return zero, ErrComponentDoNotExists
	default:
		component := c.components[componentIndex]
		return component, nil
	}
}

func (c *componentsArray[Component]) GetEntities() []EntityID { return c.componentsEntities }
func (c *componentsArray[Component]) GetAnyComponent(entity EntityID) (any, error) {
	return c.GetComponent(entity)
}

func (c *componentsArray[Component]) RemoveComponent(entity EntityID) {
	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return
	}
	componentIndex := c.entitiesComponents[entityIndex]
	if componentIndex == noEntity || componentIndex == noComponent {
		return
	}

	c.entitiesComponents[entityIndex] = noComponent

	if len(c.components)-1 == int(componentIndex) {
		c.components = c.components[:len(c.components)-1]
		c.componentsEntities = c.componentsEntities[:len(c.componentsEntities)-1]
	} else {
		movedComponent := c.components[len(c.components)-1]
		movedComponentEntity := c.componentsEntities[len(c.componentsEntities)-1]

		c.components[componentIndex] = movedComponent
		c.componentsEntities[componentIndex] = movedComponentEntity

		c.components = c.components[:len(c.components)-1]
		c.componentsEntities = c.componentsEntities[:len(c.componentsEntities)-1]

		movedEntityIndex := movedComponentEntity.Index()
		c.entitiesComponents[movedEntityIndex] = componentIndex
	}

	// listeners
	entities := []EntityID{entity}
	for _, listener := range c.onRemove {
		listener(entities)
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
	entityIndex := entity.Index()
	if entityIndex >= len(c.entitiesComponents) {
		return
	}
	componentIndex := c.entitiesComponents[entityIndex]
	c.entitiesComponents[entityIndex] = noEntity
	if componentIndex != noComponent && componentIndex != noEntity {
		// listeners
		entities := []EntityID{entity}
		for _, listener := range c.onRemove {
			listener(entities)
		}
	}
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
