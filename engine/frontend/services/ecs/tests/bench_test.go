package ecs_test

import (
	"frontend/services/ecs"
	"testing"
)

func BenchmarkGetComponentType(b *testing.B) {
	type Component struct{}
	for i := 0; i < b.N; i++ {
		ecs.GetComponentType(Component{})
	}
}

func BenchmarkGetComponentPointerType(b *testing.B) {
	type Component struct{}
	for i := 0; i < b.N; i++ {
		ecs.GetComponentPointerType((*Component)(nil))
	}
}

func BenchmarkComponentsArraySave(b *testing.B) {
	type Component struct{}
	world := ecs.NewComponentsArray[Component]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		world.AddEntity(entity)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		err := world.SaveComponent(entity, Component{})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkComponentsArrayGet10Times(b *testing.B) {
	type Component bool
	world := ecs.NewComponentsArray[bool]()

	for i := 0; i < b.N; i++ {
		entity := ecs.NewEntityID(uint64(i))
		world.AddEntity(entity)
		world.SaveComponent(entity, false)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			entity := ecs.NewEntityID(uint64(i))
			_, _ = world.GetComponent(entity)
		}
	}
}

func BenchmarkSaveComponent(b *testing.B) {
	type Component struct{}
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	entity := world.NewEntity()

	for i := 0; i < b.N; i++ {
		world.SaveComponent(entity, Component{})
	}
}

func BenchmarkGetComponent(b *testing.B) {
	type Component struct{}
	world := ecs.NewWorld()

	otherEntitiesPresent := 100
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	entities := make([]ecs.EntityID, b.N)
	for i := 0; i < b.N; i++ {
		entity := world.NewEntity()
		world.SaveComponent(entity, Component{})
		entities[i] = entity
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecs.GetComponent[Component](world, entities[i])
	}
}

func BenchmarkQueryEntitiesWithComponents(b *testing.B) {
	type RequiredComponent struct{}
	world := ecs.NewWorld()

	otherEntitiesPresent := 10000
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	requiredEntitiesPresent := 10000
	for i := 0; i < requiredEntitiesPresent; i++ {
		entity := world.NewEntity()
		world.SaveComponent(entity, RequiredComponent{})
	}

	componentType := ecs.GetComponentType(RequiredComponent{})
	for i := 0; i < b.N; i++ {
		world.QueryEntitiesWithComponents(componentType)
	}
}
