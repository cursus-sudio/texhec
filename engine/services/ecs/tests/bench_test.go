package ecs_test

import (
	"engine/services/ecs"
	"testing"
)

func BenchmarkGetComponent(b *testing.B) {
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for range otherEntitiesPresent {
		world.NewEntity()
	}

	entitiesCount := min(b.N, 10000)
	entities := make([]ecs.EntityID, entitiesCount)
	arr := ecs.GetComponentsArray[Component](world)
	for i := range entitiesCount {
		entity := world.NewEntity()
		arr.Set(entity, Component{})
		entities[i] = entity
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entityIndex := i % entitiesCount
		entity := entities[entityIndex]
		arr.Get(entity)
	}
}

func BenchmarkCreateComponents(b *testing.B) {
	world := ecs.NewWorld()
	arr := ecs.NewComponentsArray[Component](world)

	entities := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		entities[i] = entity
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := entities[i]
		arr.Set(entity, Component{})
	}
}

func BenchmarkUpdateComponents(b *testing.B) {
	world := ecs.NewWorld()
	arr := ecs.NewComponentsArray[Component](world)

	entities := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		entities[i] = entity
		arr.Set(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := entities[i]
		arr.Set(entity, Component{})
	}
}

func BenchmarkRemoveComponent(b *testing.B) {
	world := ecs.NewWorld()
	arr := ecs.NewComponentsArray[Component](world)

	entities := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		entities[i] = entity
		arr.Set(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := entities[i]
		arr.Remove(entity)
	}
}
func BenchmarkRemoveEntityWithComponent(b *testing.B) {
	world := ecs.NewWorld()
	arr := ecs.NewComponentsArray[Component](world)

	entities := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		entities[i] = entity
		arr.Set(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := entities[i]
		world.RemoveEntity(entity)
	}
}

func BenchmarkRemoveEntity(b *testing.B) {
	world := ecs.NewWorld()

	entities := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		entities[i] = entity
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := entities[i]
		world.RemoveEntity(entity)
	}
}
