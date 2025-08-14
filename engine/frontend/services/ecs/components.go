package ecs

import (
	"strings"
)

type queryKey string

func newQuerykey(slice []ComponentType) queryKey {
	var r strings.Builder
	// reserved 10 characters per element for better performance
	r.Grow(10 * len(slice))
	for _, e := range slice {
		r.WriteString(e.componentType.String())
		r.WriteString(",")
	}
	return queryKey(r.String())
}

//

type query struct {
	Dependencies map[ComponentType]any // this is faster []ComponentType
	Entities     map[EntityID]any      // this is faster []EntityID
	CachedRes    []EntityID
}

func newQuery(componentTypes []ComponentType, res []EntityID) *query {
	dependencies := make(map[ComponentType]any, len(componentTypes))
	for _, componentType := range componentTypes {
		dependencies[componentType] = nil
	}
	entities := make(map[EntityID]any, len(res))
	for _, entity := range res {
		entities[entity] = nil
	}
	return &query{dependencies, entities, nil}
}

func (query *query) RemoveEntity(entity EntityID) {
	delete(query.Entities, entity)
	query.CachedRes = nil
}

func (query *query) RemoveComponent(entity EntityID, componentType ComponentType) {
	_, ok := query.Dependencies[componentType]
	if !ok {
		return
	}
	query.RemoveEntity(entity)
}

func (query *query) AddEntity(entity EntityID) {
	query.Entities[entity] = nil
	query.CachedRes = append(query.CachedRes, entity)
}

func (query *query) Res() []EntityID {
	if query.CachedRes != nil {
		return query.CachedRes
	}
	res := make([]EntityID, 0, len(query.Entities))
	for entity := range query.Entities {
		res = append(res, entity)
	}
	query.CachedRes = res
	return res
}

//

type componentsImpl struct {
	entityComponents map[EntityID]map[ComponentType]*Component
	componentEntity  map[ComponentType]map[EntityID]*Component
	cachedQueries    map[queryKey]*query

	shouldDie bool
}

func newComponents() *componentsImpl {
	return &componentsImpl{
		entityComponents: make(map[EntityID]map[ComponentType]*Component),
		componentEntity:  make(map[ComponentType]map[EntityID]*Component),

		cachedQueries: make(map[queryKey]*query),
	}
}

func (components *componentsImpl) SaveComponent(entityId EntityID, component Component) error {
	componentType := GetComponentType(component)
	if components.entityComponents[entityId] == nil {
		return ErrEntityDoNotExists
	}

	entityHadComponent := components.componentEntity[componentType] != nil
	if !entityHadComponent {
		components.componentEntity[componentType] = make(map[EntityID]*Component)
	}

	components.entityComponents[entityId][componentType] = &component
	components.componentEntity[componentType][entityId] = &component

	if entityHadComponent {
		return nil
	}
	// manage cache
	for _, query := range components.cachedQueries {
		_, ok := query.Dependencies[componentType]
		if !ok {
			continue
		}
		dependenciesNeeded := len(query.Dependencies)
	entityDependencies:
		for componentType := range components.entityComponents[entityId] {
			if _, ok := query.Dependencies[componentType]; !ok {
				continue
			}
			dependenciesNeeded--
			if dependenciesNeeded == 0 {
				query.AddEntity(entityId)
				break entityDependencies
			}
		}
	}

	return nil
}

func (components *componentsImpl) GetComponent(entityId EntityID, componentType ComponentType) (Component, error) {
	entity, ok := components.entityComponents[entityId]
	if !ok {
		return nil, ErrEntityDoNotExists
	}
	componentPtr, ok := entity[componentType]
	if !ok {
		return nil, ErrComponentDoNotExists
	}
	return *componentPtr, nil
}

func (components *componentsImpl) RemoveComponent(entityId EntityID, componentType ComponentType) {
	delete(components.entityComponents[entityId], componentType)
	delete(components.componentEntity[componentType], entityId)
	// manage cache
	for _, query := range components.cachedQueries {
		query.RemoveComponent(entityId, componentType)
	}
}

func (components *componentsImpl) AddEntity(entity EntityID) {
	components.entityComponents[entity] = make(map[ComponentType]*Component)
}

func (components *componentsImpl) RemoveEntity(entityID EntityID) {
	entityComponents, ok := components.entityComponents[entityID]
	if !ok {
		return
	}
	for componentType := range entityComponents {
		delete(components.componentEntity[componentType], entityID)
	}
	delete(components.entityComponents, entityID)
	for _, query := range components.cachedQueries {
		query.RemoveEntity(entityID)
	}
}

func (components *componentsImpl) GetEntitiesWithComponents(componentTypes ...ComponentType) []EntityID {
	if len(componentTypes) == 0 {
		return nil
	}

	queryKey := newQuerykey(componentTypes)
	if res, ok := components.cachedQueries[queryKey]; ok {
		return res.Res()
	}

	firstType := componentTypes[0]
	entitiesMap, ok := components.componentEntity[firstType]
	if !ok || len(entitiesMap) == 0 {
		return nil
	}

	intersection := make(map[EntityID]struct{}, len(entitiesMap))
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

	entitiesSlice := make([]EntityID, 0, len(intersection))
	for entity := range intersection {
		entitiesSlice = append(entitiesSlice, entity)
	}
	components.cachedQueries[queryKey] = newQuery(componentTypes, entitiesSlice)

	return entitiesSlice
}
