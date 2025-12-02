package ecs

import (
	"engine/services/datastructures"
	"errors"
	"reflect"
)

var ErrInvalidType error = errors.New("expected an error component")

// interface

type EntityComponent[Component any] interface {
	Get() (Component, error)
	Set(Component)
	Remove()
}

type AnyComponentArray interface {
	AnyTransaction() AnyComponentsArrayTransaction

	// any component array operations
	// can return:
	// - ErrEntityDoNotExists
	SaveAnyComponent(EntityID, any) error // upsert
	// differs from save component by not triggering events
	RemoveComponent(EntityID)

	// any component array getters
	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetAnyComponent(entity EntityID) (any, error)
	GetEntities() []EntityID

	// any component array listeners
	BeforeAdd(func([]EntityID))
	BeforeChange(func([]EntityID))
	BeforeRemove(func([]EntityID))

	OnAdd(func([]EntityID))
	OnChange(func([]EntityID))
	OnRemove(func([]EntityID))
}

type ComponentsArray[Component any] interface {
	AnyComponentArray
	Transaction() ComponentsArrayTransaction[Component]

	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) error // upsert

	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetComponent(entity EntityID) (Component, error)
}

// impl

type listener uint8

const (
	addListener listener = iota
	changeListener
	removeListener
)

type componentsArray[Component any] struct {
	equal      func(Component, Component) bool
	entities   datastructures.SparseSet[EntityID]
	components datastructures.SparseArray[EntityID, Component]

	// queries are used for change and remove listeners
	queries []*liveQuery

	beforeListenersOrder []listener
	beforeAdd            []func([]EntityID)
	beforeChange         []func([]EntityID)
	beforeRemove         []func([]EntityID)

	listenersOrder []listener
	onAdd          []func([]EntityID)
	onChange       []func([]EntityID)
	onRemove       []func([]EntityID)
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

		beforeListenersOrder: make([]listener, 0),
		beforeAdd:            make([]func([]EntityID), 0),
		beforeChange:         make([]func([]EntityID), 0),
		beforeRemove:         make([]func([]EntityID), 0),

		listenersOrder: make([]listener, 0),
		onAdd:          make([]func([]EntityID), 0),
		onChange:       make([]func([]EntityID), 0),
		onRemove:       make([]func([]EntityID), 0),
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
	added := !ok
	entities := []EntityID{entity}
	if added {
		for _, listener := range c.beforeAdd {
			listener(entities)
		}
	} else {
		for _, listener := range c.beforeChange {
			listener(entities)
		}
	}
	c.components.Set(entity, component)
	if added {
		for _, listener := range c.onAdd {
			listener(entities)
		}
	} else {
		for _, listener := range c.onChange {
			listener(entities)
		}
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
	entities := []EntityID{entity}
	if _, ok := c.components.Get(entity); !ok {
		return
	}
	for _, listener := range c.beforeRemove {
		listener(entities)
	}
	c.components.Remove(entity)

	for _, listener := range c.onRemove {
		listener(entities)
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

func (c *componentsArray[Component]) BeforeAdd(listener func([]EntityID)) {
	c.beforeListenersOrder = append(c.beforeListenersOrder, addListener)
	c.beforeAdd = append(c.beforeAdd, listener)
}
func (c *componentsArray[Component]) BeforeChange(listener func([]EntityID)) {
	c.beforeListenersOrder = append(c.beforeListenersOrder, changeListener)
	c.beforeChange = append(c.beforeChange, listener)
}
func (c *componentsArray[Component]) BeforeRemove(listener func([]EntityID)) {
	c.beforeListenersOrder = append(c.beforeListenersOrder, removeListener)
	c.beforeRemove = append(c.beforeRemove, listener)
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
