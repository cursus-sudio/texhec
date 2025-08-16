package ecs

import (
	"sort"
	"strings"
)

type queryKey string

func newQueryKey(components []ComponentType) queryKey {
	resultLen := 0
	elements := make([]string, len(components))
	for i, component := range components {
		element := component.componentType.String()
		resultLen += len(element) + 1
		elements[i] = element
	}
	sort.Strings(elements)
	builder := strings.Builder{}
	builder.Grow(resultLen)
	for _, element := range elements {
		builder.WriteString(element)
		builder.WriteString(",")
	}
	return queryKey(builder.String())
}

//

type liveQuery struct {
	dependencies   map[ComponentType]any // this is faster []ComponentType
	entities       map[EntityID]any      // this is faster []EntityID
	cachedEntities []EntityID
	onRemove       []func([]EntityID)
	onAdd          []func([]EntityID)
}

func newLiveQuery(
	componentTypes []ComponentType,
	res []EntityID,
) *liveQuery {
	dependencies := make(map[ComponentType]any, len(componentTypes))
	for _, componentType := range componentTypes {
		dependencies[componentType] = nil
	}
	entities := make(map[EntityID]any, len(res))
	for _, entity := range res {
		entities[entity] = nil
	}
	return &liveQuery{
		dependencies: dependencies,
		entities:     entities,
	}
}

func (query *liveQuery) OnAdd(listener func([]EntityID)) {
	listener(query.Entities())
	query.onAdd = append(query.onAdd, listener)
}

func (query *liveQuery) OnRemove(listener func([]EntityID)) {
	query.onRemove = append(query.onRemove, listener)
}

func (query *liveQuery) Entities() []EntityID {
	if query.cachedEntities == nil {
		entities := make([]EntityID, 0, len(query.entities))
		for entity := range query.entities {
			entities = append(entities, entity)
		}
		query.cachedEntities = entities
	}
	return query.cachedEntities
}

//

func (query *liveQuery) RemoveEntity(entity EntityID) {
	_, ok := query.entities[entity]
	if !ok {
		return
	}
	delete(query.entities, entity)
	query.cachedEntities = nil
	rmArgs := []EntityID{entity}
	for _, listener := range query.onRemove {
		listener(rmArgs)
	}
}

func (query *liveQuery) RemoveComponent(entity EntityID, componentType ComponentType) {
	_, ok := query.dependencies[componentType]
	if !ok {
		return
	}
	query.RemoveEntity(entity)
}

func (query *liveQuery) AddEntities(entities []EntityID) {
	for _, entity := range entities {
		query.entities[entity] = nil
	}
	query.cachedEntities = append(query.cachedEntities, entities...)
	for _, listener := range query.onAdd {
		listener(entities)
	}
}

//

type componentsImpl struct {
	entityComponents map[EntityID]map[ComponentType]*Component
	componentEntity  map[ComponentType]map[EntityID]*Component
	cachedQueries    map[queryKey]*liveQuery

	shouldDie bool
}

func newComponents() *componentsImpl {
	return &componentsImpl{
		entityComponents: make(map[EntityID]map[ComponentType]*Component),
		componentEntity:  make(map[ComponentType]map[EntityID]*Component),

		cachedQueries: make(map[queryKey]*liveQuery, 0),
	}
}

func (components *componentsImpl) SaveComponent(entityId EntityID, component Component) error {
	componentType := GetComponentType(component)
	if components.entityComponents[entityId] == nil {
		return ErrEntityDoNotExists
	}

	_, entityHadComponent := components.entityComponents[entityId][componentType]
	if components.componentEntity[componentType] == nil {
		components.componentEntity[componentType] = make(map[EntityID]*Component)
	}

	components.entityComponents[entityId][componentType] = &component
	components.componentEntity[componentType][entityId] = &component

	if entityHadComponent {
		return nil
	}
	// manage cache
	for _, query := range components.cachedQueries {
		_, ok := query.dependencies[componentType]
		if !ok {
			continue
		}
		dependenciesNeeded := len(query.dependencies)
		entityComponents := components.entityComponents[entityId]
		for componentType := range entityComponents {
			if _, ok := query.dependencies[componentType]; !ok {
				continue
			}
			dependenciesNeeded--
			if dependenciesNeeded == 0 {
				break
			}
		}
		if dependenciesNeeded == 0 {
			query.AddEntities([]EntityID{entityId})
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

	return entitiesSlice
}

func (components *componentsImpl) GetEntitiesWithComponentsQuery(componentTypes ...ComponentType) LiveQuery {
	key := newQueryKey(componentTypes)
	if query, ok := components.cachedQueries[key]; ok {
		return query
	}
	entities := components.GetEntitiesWithComponents(componentTypes...)
	queryRes := newLiveQuery(componentTypes, entities)
	components.cachedQueries[key] = queryRes
	return queryRes
}
