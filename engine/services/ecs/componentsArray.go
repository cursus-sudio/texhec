package ecs

import (
	"engine/services/datastructures"
	"errors"
	"reflect"
)

var ErrInvalidType error = errors.New("expected an error component")

// interface

type BeforeGet func()

type AnyComponentArray interface {
	GetAny(entity EntityID) (any, bool)
	GetEntities() []EntityID

	// when type doesn't match error is returned
	SetAny(EntityID, any) error
	Remove(EntityID)

	// configuration
	// on dependency change its also applied here
	AddDependency(AnyComponentArray)
	AddDirtySet(DirtySet)
	BeforeGet(BeforeGet)
}

type ComponentsArray[Component any] interface {
	AnyComponentArray
	Get(entity EntityID) (Component, bool)

	Set(EntityID, Component)

	// configuration
	SetEmpty(Component)
}

// impl

type componentsArray[Component any] struct {
	entities   entitiesInterface
	equal      func(Component, Component) bool
	empty      Component
	components datastructures.SparseArray[EntityID, Component]

	dependencies []AnyComponentArray
	dirtySets    datastructures.Set[DirtySet]
	beforeGets   []BeforeGet
}

func NewComponentsArray[Component any](entities entitiesInterface) *componentsArray[Component] {
	equal := func(Component, Component) bool { return false }
	if reflect.TypeFor[Component]().Comparable() {
		equal = func(c1, c2 Component) bool { return any(c1) == any(c2) }
	}
	array := &componentsArray[Component]{
		entities: entities,
		equal:    equal,
		// empty: default,
		components: datastructures.NewSparseArray[EntityID, Component](),

		dependencies: nil,
		dirtySets:    datastructures.NewSet[DirtySet](),
		beforeGets:   nil,
	}
	return array
}

func (c *componentsArray[Component]) Set(entity EntityID, component Component) {
	value, ok := c.components.Get(entity)
	if ok && c.equal(value, component) {
		return
	}
	c.entities.EnsureExists(entity)
	c.components.Set(entity, component)
	for _, dirtyFlags := range c.dirtySets.Get() {
		dirtyFlags.Dirty(entity)
	}
}

func (c *componentsArray[Component]) SetAny(entity EntityID, anyComponent any) error {
	component, ok := anyComponent.(Component)
	if !ok {
		return ErrInvalidType
	}
	c.Set(entity, component)
	return nil
}

func (c *componentsArray[Component]) SetEmpty(empty Component) {
	c.empty = empty
}

func (c *componentsArray[Component]) Remove(entity EntityID) {
	entities := []EntityID{entity}
	if _, ok := c.components.Get(entity); !ok {
		return
	}
	c.components.Remove(entity)
	for _, dirtyFlag := range c.dirtySets.Get() {
		for _, entity := range entities {
			dirtyFlag.Dirty(entity)
		}
	}
}

func (c *componentsArray[Component]) Get(entity EntityID) (Component, bool) {
	for _, beforeGet := range c.beforeGets {
		beforeGet()
	}
	if value, ok := c.components.Get(entity); !ok {
		return c.empty, false
	} else {
		return value, true
	}
}

func (c *componentsArray[Component]) GetEntities() []EntityID {
	for _, beforeGet := range c.beforeGets {
		beforeGet()
	}
	return c.components.GetIndices()
}

func (c *componentsArray[Component]) GetAny(entity EntityID) (any, bool) {
	for _, beforeGet := range c.beforeGets {
		beforeGet()
	}
	return c.Get(entity)
}

//

func (c *componentsArray[Component]) AddDependency(dependency AnyComponentArray) {
	c.dependencies = append(c.dependencies, dependency)
	for _, dirtySet := range c.dirtySets.Get() {
		dependency.AddDirtySet(dirtySet)
	}
}
func (c *componentsArray[Component]) AddDirtySet(dirtySet DirtySet) {
	if _, ok := c.dirtySets.GetIndex(dirtySet); ok {
		return
	}
	for _, entity := range c.GetEntities() {
		dirtySet.Dirty(entity)
	}
	for _, dependency := range c.dependencies {
		dependency.AddDirtySet(dirtySet)
	}
	c.dirtySets.Add(dirtySet)
}
func (c *componentsArray[Component]) BeforeGet(listener BeforeGet) {
	c.beforeGets = append(c.beforeGets, listener)
}
