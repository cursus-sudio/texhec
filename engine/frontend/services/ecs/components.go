package ecs

import "sync"

//

type componentsImpl struct {
	mutex sync.Locker

	componentArrays  map[ComponentType]any // any is *componentsArray[ComponentType]
	entityComponents map[EntityID]map[ComponentType]*Component
	componentEntity  map[ComponentType]map[EntityID]*Component

	cachedQueries    map[queryKey]*liveQuery
	dependentQueries map[ComponentType][]*liveQuery
}

func newComponents() *componentsImpl {
	return &componentsImpl{
		mutex:           &sync.Mutex{},
		componentArrays: make(map[ComponentType]any),

		entityComponents: make(map[EntityID]map[ComponentType]*Component),
		componentEntity:  make(map[ComponentType]map[EntityID]*Component),

		cachedQueries:    make(map[queryKey]*liveQuery, 0),
		dependentQueries: make(map[ComponentType][]*liveQuery, 0),
	}
}

// func (components *componentsImpl) getComponentArray(componentType ComponentType) any {
// 	if arr, ok := components.componentArrays[componentType]; ok {
// 		return arr
// 	}
// 	// newArray := newComponentsArray
// 	return nil
// }

func (components *componentsImpl) SaveComponent(entityID EntityID, component Component) error {
	componentType := GetComponentType(component)
	if components.entityComponents[entityID] == nil {
		return ErrEntityDoNotExists
	}

	_, entityHadComponent := components.entityComponents[entityID][componentType]
	if components.componentEntity[componentType] == nil {
		components.componentEntity[componentType] = make(map[EntityID]*Component)
	}

	components.entityComponents[entityID][componentType] = &component
	components.componentEntity[componentType][entityID] = &component

	dependentQueries, _ := components.dependentQueries[componentType]
	if entityHadComponent {
		for _, query := range dependentQueries {
			if _, ok := query.entities.GetIndex(entityID); !ok {
				continue
			}
			query.Changed([]EntityID{entityID})
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
	dependentQueries, _ := components.dependentQueries[componentType]
	for _, query := range dependentQueries {
		query.RemoveEntity(entityId)
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

func (components *componentsImpl) QueryEntitiesWithComponents(componentTypes ...ComponentType) LiveQuery {
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

func (components *componentsImpl) LockComponents()   { components.mutex.Lock() }
func (components *componentsImpl) UnlockComponents() { components.mutex.Unlock() }
