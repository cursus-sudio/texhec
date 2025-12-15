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
	GetAnyComponent(entity EntityID) (any, bool)
	GetEntities() []EntityID

	// when type doesn't match error is returned
	SaveAnyComponent(EntityID, any) error
	RemoveComponent(EntityID)

	// on dependency change its also applied here
	AddDependency(AnyComponentArray)
	AddDirtySet(DirtySet)
	BeforeGet(BeforeGet)
}

type ComponentsArray[Component any] interface {
	AnyComponentArray
	GetComponent(entity EntityID) (Component, bool)

	SaveComponent(EntityID, Component)
}

// impl

type componentsArray[Component any] struct {
	equal      func(Component, Component) bool
	components datastructures.SparseArray[EntityID, Component]

	dependencies []AnyComponentArray
	dirtySets    datastructures.Set[DirtySet]
	beforeGets   []BeforeGet
}

func NewComponentsArray[Component any](entities datastructures.SparseSet[EntityID]) *componentsArray[Component] {
	equal := func(Component, Component) bool { return false }
	if reflect.TypeFor[Component]().Comparable() {
		equal = func(c1, c2 Component) bool { return any(c1) == any(c2) }
	}
	array := &componentsArray[Component]{
		equal:      equal,
		components: datastructures.NewSparseArray[EntityID, Component](),

		dirtySets: datastructures.NewSet[DirtySet](),
	}
	return array
}

func (c *componentsArray[Component]) SaveComponent(entity EntityID, component Component) {
	value, ok := c.components.Get(entity)
	if ok && c.equal(value, component) {
		return
	}
	entities := []EntityID{entity}
	c.components.Set(entity, component)
	for _, dirtyFlags := range c.dirtySets.Get() {
		for _, entity := range entities {
			dirtyFlags.Dirty(entity)
		}
	}
	return
}

func (c *componentsArray[Component]) SaveAnyComponent(entity EntityID, anyComponent any) error {
	component, ok := anyComponent.(Component)
	if !ok {
		return ErrInvalidType
	}
	c.SaveComponent(entity, component)
	return nil
}

func (c *componentsArray[Component]) RemoveComponent(entity EntityID) {
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

func (c *componentsArray[Component]) GetComponent(entity EntityID) (Component, bool) {
	for _, beforeGet := range c.beforeGets {
		beforeGet()
	}
	var zero Component
	if value, ok := c.components.Get(entity); !ok {
		return zero, false
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

func (c *componentsArray[Component]) GetAnyComponent(entity EntityID) (any, bool) {
	for _, beforeGet := range c.beforeGets {
		beforeGet()
	}
	return c.GetComponent(entity)
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
