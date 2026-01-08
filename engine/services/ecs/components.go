package ecs

import (
	"engine/services/datastructures"
	"errors"
	"reflect"
)

// interface

type componentType struct {
	componentType reflect.Type
}

func (t *componentType) String() string { return t.componentType.String() }

func newComponentType(t reflect.Type) componentType {
	return componentType{componentType: t}
}

//

type Component interface{}

func getComponentType(component Component) componentType {
	typeOfComponent := reflect.TypeOf(component)
	kind := typeOfComponent.Kind()
	if kind != reflect.Struct && kind != reflect.Array {
		panic("component has to be a struct or array (cannot use pointers under the hood)")
	}
	return newComponentType(typeOfComponent)
}

//
//
//

var (
	ErrComponentDoNotExists error = errors.New("component do not exists")
	ErrEntityDoNotExists    error = errors.New("entity do not exists")
)

type componentsInterface interface {
	// any is ComponentArray[ComponentType]
	// GetArray(ComponentType) (any, error)
	Components() ComponentsStorage

	// returns for with all listed component types
	// the same live query should be returned for the same input
}

// impl

type componentsImpl struct {
	storage *componentsStorage
}

func (components *componentsImpl) Components() ComponentsStorage { return components.storage }

func (components *componentsImpl) RemoveEntity(entity EntityID) {
	for _, arr := range components.storage.arrays {
		arr.Remove(entity)
	}
}

func newComponents(entities datastructures.SparseSet[EntityID]) *componentsImpl {
	return &componentsImpl{
		storage: newComponentsStorage(entities),
	}
}

//

type arraysSharedInterface interface {
	AnyComponentArray
	// this adds listeners for change and remove
}

type componentsStorage struct {
	arrays              map[componentType]arraysSharedInterface // any is *componentsArray[ComponentType]
	entities            datastructures.SparseSet[EntityID]
	onArrayAddListeners map[componentType][]func(arraysSharedInterface)
}

type ComponentsStorage *componentsStorage

func newComponentsStorage(entities datastructures.SparseSet[EntityID]) ComponentsStorage {
	return &componentsStorage{
		arrays:              make(map[componentType]arraysSharedInterface),
		entities:            entities,
		onArrayAddListeners: make(map[componentType][]func(arraysSharedInterface)),
	}
}

func GetComponentsArray[Component any](world World) ComponentsArray[Component] {
	components := world.Components()
	var zero Component
	componentType := getComponentType(zero)

	if array, ok := components.arrays[componentType]; ok {
		return array.(ComponentsArray[Component])
	}
	array := NewComponentsArray[Component](world)
	components.arrays[componentType] = array
	//
	listeners := components.onArrayAddListeners[componentType]
	for _, listener := range listeners {
		listener(array)
	}
	delete(components.onArrayAddListeners, componentType)
	return array
}

func SaveComponent[Component any](
	w World,
	entity EntityID,
	component Component,
) {
	GetComponentsArray[Component](w).Set(entity, component)
}

func GetComponent[Component any](
	w World,
	entity EntityID,
) (Component, bool) {
	return GetComponentsArray[Component](w).
		Get(entity)
}

func RemoveComponent[Component any](
	w World,
	entity EntityID,
) {
	GetComponentsArray[Component](w).
		Remove(entity)
}

func GetEntitiesWithComponents(
	components ComponentsStorage,
	componentTypes ...componentType,
) []EntityID {
	if len(componentTypes) == 0 {
		return nil
	}

	var arrays []arraysSharedInterface
	for _, componentType := range componentTypes {
		array, ok := components.arrays[componentType]
		if !ok {
			return nil
		}
		arrays = append(arrays, array)
	}

	arrayEntities := arrays[0].GetEntities()
	arrays = arrays[1:]
	finalEntities := []EntityID{}
arrayEntities:
	for _, entity := range arrayEntities {
		for _, array := range arrays {
			if _, ok := array.GetAny(entity); !ok {
				continue arrayEntities
			}
		}
		finalEntities = append(finalEntities, entity)
	}
	return finalEntities
}
