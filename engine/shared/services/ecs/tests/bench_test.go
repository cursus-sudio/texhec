package ecs_test

import (
	"shared/services/datastructures"
	"shared/services/ecs"
	"testing"
)

func BenchmarkGetComponentType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ecs.GetComponentType(Component{})
	}
}

func BenchmarkGetComponentPointerType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ecs.GetComponentPointerType((*Component)(nil))
	}
}

func BenchmarkSaveComponentInWorld(b *testing.B) {
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	entity := world.NewEntity()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ecs.SaveComponent(world.Components(), entity, Component{})
	}
}

func BenchmarkGetComponentInWorld(b *testing.B) {
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	entities := make([]ecs.EntityID, b.N)
	arr := ecs.GetComponentsArray[Component](world.Components())
	transaction := arr.Transaction()
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		transaction.SaveComponent(entity, Component{})
		entities[i] = entity
	}
	ecs.FlushMany(transaction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecs.GetComponent[Component](world.Components(), entities[i])
	}
}

var globalResult int

func BenchmarkGetComponentInArray(b *testing.B) {
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	entitiesCount := min(b.N, 10000)
	entities := make([]ecs.EntityID, entitiesCount)
	arr := ecs.GetComponentsArray[Component](world.Components())
	transaction := arr.Transaction()
	for i := 0; i < entitiesCount; i++ {
		entity := world.NewEntity()
		transaction.SaveComponent(entity, Component{})
		entities[i] = entity
	}
	ecs.FlushMany(transaction)

	b.ResetTimer()
	sum := 0
	for i := 0; i < b.N; i++ {
		entityIndex := i % entitiesCount
		sum += entityIndex
		entity := entities[entityIndex]
		arr.GetComponent(entity)
	}
	globalResult = sum
}

func BenchmarkModuloToSubtractFromGetComponentInArray(b *testing.B) {
	maxI := min(b.N, 10000)

	b.ResetTimer()
	sum := 0
	for i := 0; i < b.N; i++ {
		entityIndex := i % maxI
		sum += entityIndex
	}
	globalResult = sum
}

func BenchmarkCreateComponentsInArray(b *testing.B) {
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

func BenchmarkUpdateComponentsInArray(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)
	transaction := arr.Transaction()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		transaction.SaveComponent(entity, Component{})
	}

	ecs.FlushMany(transaction)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkTransactionUpdateComponentsInArray(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)
	transaction := arr.Transaction()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		transaction.SaveComponent(entity, Component{})
	}
	ecs.FlushMany(transaction)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkGetComponentsInArray10Times(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)
	transaction := arr.Transaction()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		transaction.SaveComponent(entity, Component{})
	}
	ecs.FlushMany(transaction)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			entity := ecs.NewEntityID(uint64(i))
			_, _ = arr.GetComponent(entity)
		}
	}
}

func BenchmarkRemoveComponentInArray(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)
	transaction := arr.Transaction()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		transaction.SaveComponent(entity, Component{})
	}
	ecs.FlushMany(transaction)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.RemoveComponent(entity)
	}
}
