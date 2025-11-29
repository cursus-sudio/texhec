package ecs

import (
	"engine/services/datastructures"
	"errors"
	"reflect"
)

// interface

type ComponentType struct {
	componentType reflect.Type
}

func (t *ComponentType) String() string { return t.componentType.String() }

func newComponentType(componentType reflect.Type) ComponentType {
	return ComponentType{componentType: componentType}
}

//

type Component interface{}

func GetComponentType(component Component) ComponentType {
	typeOfComponent := reflect.TypeOf(component)
	if typeOfComponent.Kind() != reflect.Struct {
		panic("component has to be a struct (cannot use pointers under the hood)")
	}
	return newComponentType(typeOfComponent)
}

func GetComponentPointerType(componentPointer any) ComponentType {
	return newComponentType(reflect.TypeOf(componentPointer).Elem())
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

	GetAnyComponent(EntityID, ComponentType) (any, error)
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
	arrays              map[ComponentType]arraysSharedInterface // any is *componentsArray[ComponentType]
	entities            datastructures.SparseSet[EntityID]
	onArrayAddListeners map[ComponentType][]func(arraysSharedInterface)

	cachedQueries    map[queryKey]*liveQuery
	dependentQueries map[ComponentType][]*liveQuery
}

type ComponentsStorage *componentsStorage

func newComponentsStorage(entities datastructures.SparseSet[EntityID]) ComponentsStorage {
	return &componentsStorage{
		arrays:              make(map[ComponentType]arraysSharedInterface),
		entities:            entities,
		onArrayAddListeners: make(map[ComponentType][]func(arraysSharedInterface)),

		cachedQueries:    make(map[queryKey]*liveQuery, 0),
		dependentQueries: make(map[ComponentType][]*liveQuery, 0),
	}
}

func (components *componentsStorage) whenArrExists(t ComponentType, l func(arraysSharedInterface)) {
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
	componentType := GetComponentType(zero)

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
	componentTypes ...ComponentType,
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
func (i *componentsImpl) GetAnyComponent(entity EntityID, componentType ComponentType) (any, error) {
	if array, ok := i.storage.arrays[componentType]; ok {
		return array.GetAnyComponent(entity)
	}
	return nil, ErrComponentDoNotExists
}
