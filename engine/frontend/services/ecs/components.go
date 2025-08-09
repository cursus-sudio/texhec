package ecs

import (
	"reflect"
)

type componentsImpl struct {
	entityComponents map[EntityId]map[ComponentType]*Component
	componentEntity  map[ComponentType]map[EntityId]*Component
}

func newComponents() *componentsImpl {
	return &componentsImpl{
		entityComponents: make(map[EntityId]map[ComponentType]*Component),
		componentEntity:  make(map[ComponentType]map[EntityId]*Component),
	}
}

func (components *componentsImpl) ChangeTo(component any, newComponent Component) {
	componentValue := reflect.ValueOf(component)
	newComponentValue := reflect.ValueOf(newComponent)
	if componentValue.Kind() != reflect.Ptr {
		panic("component have to be pointer\n")
	}
	componentValue.Elem().Set(newComponentValue)
}

func (components *componentsImpl) SaveComponent(entityId EntityId, component Component) error {
	componentType := GetComponentType(component)

	if components.entityComponents[entityId] == nil {
		return ErrEntityDoNotExists
	}

	if components.componentEntity[componentType] == nil {
		components.componentEntity[componentType] = make(map[EntityId]*Component)
	}

	components.entityComponents[entityId][componentType] = &component
	components.componentEntity[componentType][entityId] = &component
	return nil
}

func (components *componentsImpl) GetComponents(entityId EntityId, componentsPointers ...any) error {
	for _, componentPointer := range componentsPointers {
		componentType := GetComponentPointerType(componentPointer)
		entity, ok := components.entityComponents[entityId]
		if !ok {
			return ErrEntityDoNotExists
		}
		componentValue, ok := entity[componentType]
		if !ok {
			return ErrComponentDoNotExists
		}
		components.ChangeTo(componentPointer, *componentValue)
	}
	return nil
}

func (components *componentsImpl) RemoveComponent(entityId EntityId, componentType ComponentType) {
	delete(components.entityComponents[entityId], componentType)
	delete(components.componentEntity[componentType], entityId)
}

func (components *componentsImpl) GetEntitiesWithComponents(componentTypes ...ComponentType) []EntityId {
	if len(componentTypes) == 0 {
		return nil
	}

	firstType := componentTypes[0]
	entitiesMap, ok := components.componentEntity[firstType]
	if !ok || len(entitiesMap) == 0 {
		return nil
	}

	intersection := make(map[EntityId]struct{}, len(entitiesMap))
	for entity := range entitiesMap {
		intersection[entity] = struct{}{}
	}

	for _, componentType := range componentTypes[1:] {
		entitiesMap, ok := components.componentEntity[componentType]
		if !ok || len(entitiesMap) == 0 {
			return nil
		}

		for entity := range intersection {
			if _, exists := entitiesMap[entity]; !exists {
				delete(intersection, entity)
			}
		}

		if len(intersection) == 0 {
			return nil
		}
	}

	entitiesSlice := make([]EntityId, 0, len(intersection))
	for entity := range intersection {
		entitiesSlice = append(entitiesSlice, entity)
	}

	return entitiesSlice
}
