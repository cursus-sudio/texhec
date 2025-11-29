package ecs

import (
	"engine/services/datastructures"
	"errors"
	"reflect"
)

var ErrInvalidType error = errors.New("expected an error component")

// interface

type AnyComponentArray interface {
	AnyTransaction() AnyComponentsArrayTransaction

	// can return:
	// - ErrEntityDoNotExists
	SaveAnyComponent(EntityID, any) error // upsert
	// differs from save component by not triggering events
	RemoveComponent(EntityID)

	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetAnyComponent(entity EntityID) (any, error)
	GetEntities() []EntityID

	OnAdd(func([]EntityID))
	OnChange(func([]EntityID))
	OnRemove(func([]EntityID))
}

type EntityComponent[Component any] interface {
	Get() (Component, error)
	Set(Component)
	Remove()
}

type ComponentsArray[Component any] interface {
	Transaction() ComponentsArrayTransaction[Component]
	AnyTransaction() AnyComponentsArrayTransaction

	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) error // upsert
	SaveAnyComponent(EntityID, any) error    // upsert

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
	OnRemoveComponents(func([]EntityID, []Component))
}

// impl

type listener uint8

const (
	addListener listener = iota
	changeListener
	removeListener
	removeComponentsListener
)

type componentsArray[Component any] struct {
	equal      func(Component, Component) bool
	entities   datastructures.SparseSet[EntityID]
	components datastructures.SparseArray[EntityID, Component]

	// queries are used for change and remove listeners
	queries []*liveQuery

	listenersOrder     []listener
	onAdd              []func([]EntityID)
	onChange           []func([]EntityID)
	onRemove           []func([]EntityID)
	onRemoveComponents []func([]EntityID, []Component)
}

func NewComponentsArray[Component any](entities datastructures.SparseSet[EntityID]) *componentsArray[Component] {
	equal := func(Component, Component) bool { return false }
	if reflect.TypeFor[Component]().Comparable() {
		equal = func(c1, c2 Component) bool { return any(c1) == any(c2) }
	}
	array := &componentsArray[Component]{
		equal:      equal,
		entities:   entities,
		components: datastructures.NewSparseArray[EntityID, Component](),

		listenersOrder:     make([]listener, 0),
		onAdd:              make([]func([]EntityID), 0),
		onChange:           make([]func([]EntityID), 0),
		onRemove:           make([]func([]EntityID), 0),
		onRemoveComponents: make([]func([]EntityID, []Component), 0),
	}
	return array
}

func (c *componentsArray[Component]) Transaction() ComponentsArrayTransaction[Component] {
	return newComponentsArrayTransaction(c)
}

func (c *componentsArray[Component]) AnyTransaction() AnyComponentsArrayTransaction {
	return c.Transaction()
}

func (c *componentsArray[Component]) addQueries(queries []*liveQuery) {
	c.queries = append(c.queries, queries...)
}

func (c *componentsArray[Component]) SaveComponent(entity EntityID, component Component) error {
	if ok := c.entities.Get(entity); !ok {
		return ErrEntityDoNotExists
	}
	value, ok := c.components.Get(entity)
	if ok && c.equal(value, component) {
		return nil
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

func (c *componentsArray[Component]) SaveAnyComponent(entity EntityID, anyComponent any) error {
	component, ok := anyComponent.(Component)
	if !ok {
		return ErrInvalidType
	}
	return c.SaveComponent(entity, component)
}

func (c *componentsArray[Component]) RemoveComponent(entity EntityID) {
	component, _ := c.components.Get(entity)
	if removed := c.components.Remove(entity); !removed {
		return
	}
	entities := []EntityID{entity}
	components := []Component{component}

	removeI := 0
	removeComponentsI := 0
	for _, listener := range c.listenersOrder {
		switch listener {
		case addListener:
		case changeListener:
		case removeListener:
			if len(entities) != 0 {
				c.onRemove[removeI](entities)
			}
			removeI++
		case removeComponentsListener:
			if len(entities) != 0 {
				c.onRemoveComponents[removeComponentsI](entities, components)
			}
			removeComponentsI++
		}
	}
}

type entityComponent[Component any] struct {
	get func() (Component, error)
	set func(Component)
	del func()
}

func newEntityComponent[Component any](
	entity EntityID,
	get func(EntityID) (Component, error),
	set func(EntityID, Component),
	del func(EntityID),
) EntityComponent[Component] {
	return entityComponent[Component]{
		func() (Component, error) { return get(entity) },
		func(c Component) { set(entity, c) },
		func() { del(entity) },
	}
}

func NewEntityComponent[Component any](
	get func() (Component, error),
	set func(Component),
	del func(),
) EntityComponent[Component] {
	return entityComponent[Component]{
		func() (Component, error) { return get() },
		func(c Component) { set(c) },
		func() { del() },
	}
}

func (e entityComponent[Component]) Get() (Component, error) { return e.get() }
func (e entityComponent[Component]) Set(c Component)         { e.set(c) }
func (e entityComponent[Component]) Remove()                 { e.del() }

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
	listener(c.GetEntities())
	c.listenersOrder = append(c.listenersOrder, addListener)
	c.onAdd = append(c.onAdd, listener)
}

func (c *componentsArray[Component]) OnChange(listener func([]EntityID)) {
	c.listenersOrder = append(c.listenersOrder, changeListener)
	c.onChange = append(c.onChange, listener)
}

func (c *componentsArray[Component]) OnRemove(listener func([]EntityID)) {
	c.listenersOrder = append(c.listenersOrder, removeListener)
	c.onRemove = append(c.onRemove, listener)
}

func (c *componentsArray[Component]) OnRemoveComponents(listener func([]EntityID, []Component)) {
	c.listenersOrder = append(c.listenersOrder, removeComponentsListener)
	c.onRemoveComponents = append(c.onRemoveComponents, listener)
}
