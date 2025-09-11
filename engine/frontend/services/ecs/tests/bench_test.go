package ecs_test

import (
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
	componentArray := ecs.GetComponentArray[Component](world.Components())
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		componentArray.SaveComponent(entity, Component{})
		entities[i] = entity
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecs.GetComponent[Component](world.Components(), entities[i])
	}
}

func BenchmarkCreateComponentsInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		_ = arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkDirtyCreateComponentsInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		_ = arr.DirtySaveComponent(entity, Component{})
	}
}

func BenchmarkUpdateComponentsInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
		_ = arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		_ = arr.SaveComponent(entity, Component{})
	}
}

func BenchmarkDirtyUpdateComponentsInArray10Times(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
		_ = arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		for j := 0; j < 10; j++ {
			_ = arr.DirtySaveComponent(entity, Component{})
		}
	}
}

func BenchmarkGetComponentsInArray10Times(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
		arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			entity := ecs.NewEntityID(uint64(i))
			_, _ = arr.GetComponent(entity)
		}
	}
}

func BenchmarkRemoveComponentInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
		arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.RemoveComponent(entity)
	}
}

func BenchmarkAddEntityInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
	}
}

func BenchmarkRemoveEntityWithComponentInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
		arr.SaveComponent(entity, Component{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.RemoveEntity(entity)
	}
}

func BenchmarkRemoveEntityWithoutComponentInArray(b *testing.B) {
	arr := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.AddEntity(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		arr.RemoveEntity(entity)
	}
}
