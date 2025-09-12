package ecs_test

import (
	"frontend/services/datastructures"
	"frontend/services/ecs"
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
	transaction.Flush()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecs.GetComponent[Component](world.Components(), entities[i])
	}
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

func BenchmarkDirtyCreateComponentsInArray(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.DirtySaveComponent(entity, Component{})
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

	transaction.Flush()

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
	transaction.Flush()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkDirtyUpdateComponentsInArray10Times(b *testing.B) {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	arr := ecs.NewComponentsArray[Component](entities)
	transaction := arr.Transaction()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		entities.Add(entity)
		transaction.SaveComponent(entity, Component{})
	}
	transaction.Flush()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		for j := 0; j < 10; j++ {
			arr.DirtySaveComponent(entity, Component{})
		}
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
	transaction.Flush()

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
	transaction.Flush()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.RemoveComponent(entity)
	}
}
