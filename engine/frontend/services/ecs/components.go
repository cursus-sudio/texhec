package ecs

import (
	"errors"
	"frontend/services/datastructures"
	"reflect"
	"sync"
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
	QueryEntitiesWithComponents(...ComponentType) LiveQuery

	GetAnyComponent(EntityID, ComponentType) (any, error)

	LockComponents()
	UnlockComponents()
}

// impl

type componentsImpl struct {
	mutex sync.Locker

	storage ComponentsStorage
}

func (components *componentsImpl) Components() ComponentsStorage { return components.storage }
func (components *componentsImpl) LockComponents()               { components.mutex.Lock() }
func (components *componentsImpl) UnlockComponents()             { components.mutex.Unlock() }

func (components *componentsImpl) RemoveEntity(entity EntityID) {
	for _, arr := range components.storage.arrays {
		arr.RemoveComponent(entity)
	}
}

func newComponents(entities datastructures.SparseSet[EntityID]) *componentsImpl {
	return &componentsImpl{
		mutex:   &sync.Mutex{},
		storage: newComponentsStorage(entities),
	}
}

//

type arraysSharedInterface interface {
	GetEntities() []EntityID
	GetAnyComponent(entity EntityID) (any, error)
	RemoveComponent(EntityID)

	// this adds listeners for change and remove
	addQueries([]*liveQuery)

	OnAdd(listener func([]EntityID))
	OnChange(listener func([]EntityID))
	OnRemove(listener func([]EntityID))
}

type componentsStorage struct {
	arrays     map[ComponentType]arraysSharedInterface // any is *componentsArray[ComponentType]
	entities   datastructures.SparseSet[EntityID]
	onArrayAdd map[ComponentType][]*liveQuery

	cachedQueries    map[queryKey]*liveQuery
	dependentQueries map[ComponentType][]*liveQuery
}

type ComponentsStorage *componentsStorage

func newComponentsStorage(entities datastructures.SparseSet[EntityID]) ComponentsStorage {
	return &componentsStorage{
		arrays:     make(map[ComponentType]arraysSharedInterface),
		entities:   entities,
		onArrayAdd: make(map[ComponentType][]*liveQuery),

		cachedQueries:    make(map[queryKey]*liveQuery, 0),
		dependentQueries: make(map[ComponentType][]*liveQuery, 0),
	}
}

func addDependentQueriesListeners(
	components ComponentsStorage,
	componentType ComponentType,
) {
	queries, ok := components.onArrayAdd[componentType]
	if !ok {
		return
	}
	array, ok := components.arrays[componentType]
	if !ok {
		return
	}
	delete(components.onArrayAdd, componentType)

	for _, query := range queries {
		arrays := make([]arraysSharedInterface, 0, len(query.dependencies.Get()))
		missingArrays := query.dependencies.Get()
		array.OnAdd(func(ei []EntityID) {
			for _, missingArray := range missingArrays {
				array, ok := components.arrays[missingArray]
				if !ok {
					return
				}
				missingArrays = missingArrays[1:]
				arrays = append(arrays, array)
			}
			addedEntities := make([]EntityID, 0, len(ei))
		entityLoop:
			for _, entity := range ei {
				for _, array := range arrays {
					if _, err := array.GetAnyComponent(entity); err != nil {
						continue entityLoop
					}
				}
				addedEntities = append(addedEntities, entity)
			}
			if len(addedEntities) != 0 {
				query.AddedEntities(addedEntities)
			}
		})
	}
	array.addQueries(queries)
}

func GetComponentsArray[Component any](components ComponentsStorage) ComponentsArray[Component] {
	var zero Component
	componentType := GetComponentType(zero)

	if array, ok := components.arrays[componentType]; ok {
		return array.(ComponentsArray[Component])
	}
	array := NewComponentsArray[Component](components.entities)
	components.arrays[componentType] = array
	addDependentQueriesListeners(components, componentType)
	return array
}

func SaveComponent[Component any](
	components ComponentsStorage,
	entity EntityID,
	component Component,
) error {
	return GetComponentsArray[Component](components).
		SaveComponent(entity, component)
}

func GetComponent[Component any](
	components ComponentsStorage,
	entity EntityID,
) (Component, error) {
	return GetComponentsArray[Component](components).
		GetComponent(entity)
}

func RemoveComponent[Component any](
	components ComponentsStorage,
	entity EntityID,
) {
	GetComponentsArray[Component](components).
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

func (i *componentsImpl) QueryEntitiesWithComponents(componentTypes ...ComponentType) LiveQuery {
	components := i.storage
	key := newQueryKey(componentTypes)
	if query, ok := components.cachedQueries[key]; ok {
		return query
	}
	entities := GetEntitiesWithComponents(components, componentTypes...)
	query := newLiveQuery(componentTypes)
	query.AddedEntities(entities)
	components.cachedQueries[key] = query
	for _, componentType := range componentTypes {
		dependentQueries, _ := components.dependentQueries[componentType]
		dependentQueries = append(dependentQueries, query)
		components.dependentQueries[componentType] = dependentQueries

		onAdd, _ := components.onArrayAdd[componentType]
		onAdd = append(onAdd, query)
		components.onArrayAdd[componentType] = onAdd

		if _, ok := components.arrays[componentType]; ok {
			addDependentQueriesListeners(components, componentType)
		}
	}
	return query
}
func (i *componentsImpl) GetAnyComponent(entity EntityID, componentType ComponentType) (any, error) {
	if array, ok := i.storage.arrays[componentType]; ok {
		return array.GetAnyComponent(entity)
	}
	return nil, ErrComponentDoNotExists
}
