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
	if typeOfComponent.Kind() != reflect.Struct {
		panic("component has to be a struct (cannot use pointers under the hood)")
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
	Query() LiveQueryBuilder
}

// impl

type componentsImpl struct {
	storage *componentsStorage
}

func (components *componentsImpl) Components() ComponentsStorage { return components.storage }

func (components *componentsImpl) RemoveEntity(entity EntityID) {
	for _, arr := range components.storage.arrays {
		arr.RemoveComponent(entity)
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
	addQueries([]*liveQuery)
}

type componentsStorage struct {
	arrays              map[componentType]arraysSharedInterface // any is *componentsArray[ComponentType]
	entities            datastructures.SparseSet[EntityID]
	onArrayAddListeners map[componentType][]func(arraysSharedInterface)

	cachedQueries    map[queryKey]*liveQuery
	dependentQueries map[componentType][]*liveQuery
}

type ComponentsStorage *componentsStorage

func newComponentsStorage(entities datastructures.SparseSet[EntityID]) ComponentsStorage {
	return &componentsStorage{
		arrays:              make(map[componentType]arraysSharedInterface),
		entities:            entities,
		onArrayAddListeners: make(map[componentType][]func(arraysSharedInterface)),

		cachedQueries:    make(map[queryKey]*liveQuery, 0),
		dependentQueries: make(map[componentType][]*liveQuery, 0),
	}
}

func (components *componentsStorage) whenArrExists(t componentType, l func(arraysSharedInterface)) {
	if arr, ok := components.arrays[t]; ok {
		l(arr)
		return
	}
	onAdd, _ := components.onArrayAddListeners[t]
	onAdd = append(onAdd, l)
	components.onArrayAddListeners[t] = onAdd
}

func GetComponentsArray[Component any](world World) ComponentsArray[Component] {
	components := world.Components()
	var zero Component
	componentType := getComponentType(zero)

	if array, ok := components.arrays[componentType]; ok {
		return array.(ComponentsArray[Component])
	}
	array := NewComponentsArray[Component](components.entities)
	components.arrays[componentType] = array
	//
	listeners, _ := components.onArrayAddListeners[componentType]
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
) error {
	return GetComponentsArray[Component](w).
		SaveComponent(entity, component)
}

func GetComponent[Component any](
	w World,
	entity EntityID,
) (Component, error) {
	return GetComponentsArray[Component](w).
		GetComponent(entity)
}

func RemoveComponent[Component any](
	w World,
	entity EntityID,
) {
	GetComponentsArray[Component](w).
		RemoveComponent(entity)
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
			if _, err := array.GetAnyComponent(entity); err != nil {
				continue arrayEntities
			}
		}
		finalEntities = append(finalEntities, entity)
	}
	return finalEntities
}

func (i *componentsImpl) Query() LiveQueryBuilder {
	return newLiveQueryFactory(i)
}
