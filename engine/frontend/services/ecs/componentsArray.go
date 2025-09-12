package ecs

import (
	"frontend/services/datastructures"
	"math"
	"sync"
)

// interface

type ComponentsArray[Component any] interface {
	Transaction() ComponentsArrayTransaction[Component]

	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) error // upsert
	// differs from save component by not triggering events
	DirtySaveComponent(EntityID, Component) error // upsert
	RemoveComponent(EntityID)

	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetComponent(entity EntityID) (Component, error)
	GetAnyComponent(entity EntityID) (any, error)
	GetEntities() []EntityID

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
	entities   datastructures.SparseSet[EntityID]
	components datastructures.SparseArray[EntityID, Component]

	applyTransactionMutex sync.Mutex

	onAdd    []func([]EntityID)
	onChange []func([]EntityID)
	onRemove []func([]EntityID)
}

func NewComponentsArray[Component any](entities datastructures.SparseSet[EntityID]) ComponentsArray[Component] {
	return &componentsArray[Component]{
		entities:   entities,
		components: datastructures.NewSparseArray[EntityID, Component](),

		onAdd:    make([]func([]EntityID), 0),
		onChange: make([]func([]EntityID), 0),
		onRemove: make([]func([]EntityID), 0),
	}
}

func (c *componentsArray[Component]) Transaction() ComponentsArrayTransaction[Component] {
	return newComponentsArrayTransaction(c)
}

func (c *componentsArray[Component]) SaveComponent(entity EntityID, component Component) error {
	if ok := c.entities.Get(entity); !ok {
		return ErrEntityDoNotExists
	}
	added := c.components.Set(entity, component)
	entities := []EntityID{entity}
	if added {
		for _, listener := range c.onAdd {
			listener(entities)
		}
		return nil
	}
	for _, listener := range c.onChange {
		listener(entities)
	}
	return nil
}

func (c *componentsArray[Component]) DirtySaveComponent(entity EntityID, component Component) error {
	if ok := c.entities.Get(entity); !ok {
		return ErrEntityDoNotExists
	}
	c.components.Set(entity, component)
	return nil
}

func (c *componentsArray[Component]) RemoveComponent(entity EntityID) {
	if removed := c.components.Remove(entity); !removed {
		return
	}
	entities := []EntityID{entity}
	for _, listener := range c.onRemove {
		listener(entities)
	}
}

func (c *componentsArray[Component]) GetComponent(entity EntityID) (Component, error) {
	var zero Component
	if ok := c.entities.Get(entity); !ok {
		return zero, ErrEntityDoNotExists
	}
	if value, ok := c.components.Get(entity); !ok {
		return zero, ErrComponentDoNotExists
	} else {
		return value, nil
	}
}

func (c *componentsArray[Component]) GetEntities() []EntityID {
	return c.components.GetIndices()
}

func (c *componentsArray[Component]) GetAnyComponent(entity EntityID) (any, error) {
	return c.GetComponent(entity)
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
