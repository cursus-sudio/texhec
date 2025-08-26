package ecs

import (
	"frontend/services/datastructures"
	"sort"
	"strings"
	"sync"
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
	dependencies map[ComponentType]any // this is faster []ComponentType
	entities     datastructures.Set[EntityID]
	onRemove     []func([]EntityID)
	onChange     []func([]EntityID)
	onAdd        []func([]EntityID)
}

func newLiveQuery(
	componentTypes []ComponentType,
	res []EntityID,
) *liveQuery {
	dependencies := make(map[ComponentType]any, len(componentTypes))
	for _, componentType := range componentTypes {
		dependencies[componentType] = nil
	}
	entities := datastructures.NewSet[EntityID]()
	for _, entity := range res {
		entities.Add(entity)
	}
	return &liveQuery{
		dependencies: dependencies,
		entities:     entities,
	}
}

func (query *liveQuery) OnAdd(listener func([]EntityID)) {
	entities := query.Entities()
	if len(entities) != 0 {
		listener(entities)
	}
	query.onAdd = append(query.onAdd, listener)
}

func (query *liveQuery) OnChange(listener func([]EntityID)) {
	query.onChange = append(query.onChange, listener)
}

func (query *liveQuery) OnRemove(listener func([]EntityID)) {
	query.onRemove = append(query.onRemove, listener)
}

func (query *liveQuery) Entities() []EntityID {
	return query.entities.Get()
}

//

func (query *liveQuery) RemoveEntity(entity EntityID) {
	index, ok := query.entities.GetIndex(entity)
	if !ok {
		return
	}
	query.entities.Remove(index)
	rmArgs := []EntityID{entity}
	for _, listener := range query.onRemove {
		listener(rmArgs)
	}
}

func (query *liveQuery) AddEntities(entities []EntityID) {
	query.entities.Add(entities...)
	for _, listener := range query.onAdd {
		listener(entities)
	}
}

//

type componentsImpl struct {
	entityComponents map[EntityID]map[ComponentType]*Component
	componentEntity  map[ComponentType]map[EntityID]*Component

	cachedQueries    map[queryKey]*liveQuery
	dependentQueries map[ComponentType][]*liveQuery

	mutex *sync.RWMutex
}

func newComponents(mutex *sync.RWMutex) *componentsImpl {
	return &componentsImpl{
		entityComponents: make(map[EntityID]map[ComponentType]*Component),
		componentEntity:  make(map[ComponentType]map[EntityID]*Component),

		cachedQueries:    make(map[queryKey]*liveQuery, 0),
		dependentQueries: make(map[ComponentType][]*liveQuery, 0),

		mutex: mutex,
	}
}

func (components *componentsImpl) SaveComponent(entityID EntityID, component Component) error {
	components.mutex.Lock()
	componentType := GetComponentType(component)
	if components.entityComponents[entityID] == nil {
		components.mutex.Unlock()
		return ErrEntityDoNotExists
	}

	_, entityHadComponent := components.entityComponents[entityID][componentType]
	if components.componentEntity[componentType] == nil {
		components.componentEntity[componentType] = make(map[EntityID]*Component)
	}

	components.entityComponents[entityID][componentType] = &component
	components.componentEntity[componentType][entityID] = &component
	components.mutex.Unlock()

	dependentQueries, _ := components.dependentQueries[componentType]
	if entityHadComponent {
		for _, query := range dependentQueries {
			for _, listener := range query.onChange {
				listener([]EntityID{entityID})
			}
		}
		return nil
	}
	// manage cache
	for _, query := range dependentQueries {
		dependenciesNeeded := len(query.dependencies)
		entityComponents := components.entityComponents[entityID]
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
			query.AddEntities([]EntityID{entityID})
		}
	}

	return nil
}

func (components *componentsImpl) GetComponent(entityId EntityID, componentType ComponentType) (Component, error) {
	components.mutex.RLocker().Lock()
	defer components.mutex.RLocker().Unlock()
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
	components.mutex.Lock()
	delete(components.entityComponents[entityId], componentType)
	delete(components.componentEntity[componentType], entityId)
	components.mutex.Unlock()
	// manage cache
	dependentQueries, _ := components.dependentQueries[componentType]
	for _, query := range dependentQueries {
		query.RemoveEntity(entityId)
	}
}

func (components *componentsImpl) AddEntity(entity EntityID) {
	components.entityComponents[entity] = make(map[ComponentType]*Component)
	components.mutex.Unlock()
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
	components.mutex.Unlock()
	for _, query := range components.cachedQueries {
		query.RemoveEntity(entityID)
	}
}

func (components *componentsImpl) GetEntitiesWithComponents(componentTypes ...ComponentType) []EntityID {
	components.mutex.RLocker().Lock()
	defer components.mutex.RLocker().Unlock()
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

func (components *componentsImpl) QueryEntitiesWithComponents(componentTypes ...ComponentType) LiveQuery {
	components.mutex.RLocker().Lock()
	defer components.mutex.RLocker().Unlock()
	key := newQueryKey(componentTypes)
	if query, ok := components.cachedQueries[key]; ok {
		return query
	}
	entities := components.GetEntitiesWithComponents(componentTypes...)
	query := newLiveQuery(componentTypes, entities)
	components.cachedQueries[key] = query
	for _, componentType := range componentTypes {
		dependentQueries, _ := components.dependentQueries[componentType]
		dependentQueries = append(dependentQueries, query)
		components.dependentQueries[componentType] = dependentQueries
	}
	return query
}
