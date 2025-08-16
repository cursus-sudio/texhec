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

	entity := world.NewEntity()
	world.SaveComponent(entity, Component{})

	for i := 0; i < b.N; i++ {
		ecs.GetComponent[Component](world, entity)
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
