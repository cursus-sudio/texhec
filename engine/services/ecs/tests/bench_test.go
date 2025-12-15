package ecs_test

import (
	"engine/services/datastructures"
	"engine/services/ecs"
	"testing"
)

func BenchmarkGetComponent(b *testing.B) {
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	entitiesCount := min(b.N, 10000)
	entities := make([]ecs.EntityID, entitiesCount)
	arr := ecs.GetComponentsArray[Component](world)
	for i := 0; i < entitiesCount; i++ {
		entity := world.NewEntity()
		arr.SaveComponent(entity, Component{})
		entities[i] = entity
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityIndex := i % entitiesCount
		entity := entities[entityIndex]
		arr.GetComponent(entity)
	}
}

func BenchmarkCreateComponents(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkUpdateComponents(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkRemoveComponent(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.RemoveComponent(entity)
	}
}
