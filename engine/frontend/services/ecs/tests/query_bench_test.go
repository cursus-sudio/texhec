package ecs_test

import (
	"frontend/services/ecs"
	"testing"
)

func addQueries(world ecs.World, others ...ecs.ComponentType) {
	for _, other := range others {
		query := world.QueryEntitiesWithComponents(
			ecs.GetComponentType(Component{}),
			other,
		)
		query.OnAdd(func(ei []ecs.EntityID) {})
		query.OnChange(func(ei []ecs.EntityID) {})
		query.OnRemove(func(ei []ecs.EntityID) {})
	}
}

func BenchmarkQueryEntitiesWithComponents(b *testing.B) {
	world := ecs.NewWorld()

	otherEntitiesPresent := 10000
	for i := 0; i < otherEntitiesPresent; i++ {
		world.NewEntity()
	}

	requiredEntitiesPresent := 10000
	for i := 0; i < requiredEntitiesPresent; i++ {
		entity := world.NewEntity()
		ecs.SaveComponent(world.Components(), entity, Component{})
	}

	b.ResetTimer()

	componentType := ecs.GetComponentType(Component{})
	for i := 0; i < b.N; i++ {
		world.QueryEntitiesWithComponents(componentType)
	}
}

func BenchmarkSaveWithLiveQuery(b *testing.B) {
	world := ecs.NewWorld()

	type c1 struct{}

	addQueries(
		world,
		ecs.GetComponentType(c1{}),
	)
	entity := world.NewEntity()
	ecs.SaveComponent(world.Components(), entity, c1{})

	array := ecs.GetComponentArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.SaveComponent(entity, Component{})
	}
}

func BenchmarkSaveWith4LiveQueries(b *testing.B) {
	world := ecs.NewWorld()

	type c1 struct{}
	type c2 struct{}
	type c3 struct{}
	type c4 struct{}

	addQueries(
		world,
		ecs.GetComponentType(c1{}),
		ecs.GetComponentType(c2{}),
		ecs.GetComponentType(c3{}),
		ecs.GetComponentType(c4{}),
	)
	entity := world.NewEntity()
	ecs.SaveComponent(world.Components(), entity, c1{})
	ecs.SaveComponent(world.Components(), entity, c2{})
	ecs.SaveComponent(world.Components(), entity, c3{})
	ecs.SaveComponent(world.Components(), entity, c4{})

	array := ecs.GetComponentArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.SaveComponent(entity, Component{})
	}
}

func BenchmarkSaveWith7LiveQueries(b *testing.B) {
	world := ecs.NewWorld()

	type c1 struct{}
	type c2 struct{}
	type c3 struct{}
	type c4 struct{}
	type c5 struct{}
	type c6 struct{}
	type c7 struct{}

	addQueries(
		world,
		ecs.GetComponentType(c1{}),
		ecs.GetComponentType(c2{}),
		ecs.GetComponentType(c3{}),
		ecs.GetComponentType(c4{}),
		ecs.GetComponentType(c5{}),
		ecs.GetComponentType(c6{}),
		ecs.GetComponentType(c7{}),
	)
	entity := world.NewEntity()
	ecs.SaveComponent(world.Components(), entity, c1{})
	ecs.SaveComponent(world.Components(), entity, c2{})
	ecs.SaveComponent(world.Components(), entity, c3{})
	ecs.SaveComponent(world.Components(), entity, c4{})
	ecs.SaveComponent(world.Components(), entity, c5{})
	ecs.SaveComponent(world.Components(), entity, c6{})
	ecs.SaveComponent(world.Components(), entity, c7{})

	array := ecs.GetComponentArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.SaveComponent(entity, Component{})
	}
}

func Benchmark4SavesWith7LiveQueries(b *testing.B) {
	world := ecs.NewWorld()
	addQueries := func(others ...ecs.ComponentType) {
		for _, other := range others {
			query := world.QueryEntitiesWithComponents(
				ecs.GetComponentType(Component{}),
				other,
			)
			query.OnAdd(func(ei []ecs.EntityID) {})
			query.OnChange(func(ei []ecs.EntityID) {})
			query.OnRemove(func(ei []ecs.EntityID) {})
		}
	}

	type c1 struct{}
	type c2 struct{}
	type c3 struct{}
	type c4 struct{}
	type c5 struct{}
	type c6 struct{}
	type c7 struct{}

	addQueries(
		ecs.GetComponentType(c1{}),
		ecs.GetComponentType(c2{}),
		ecs.GetComponentType(c3{}),
		ecs.GetComponentType(c4{}),
		ecs.GetComponentType(c5{}),
		ecs.GetComponentType(c6{}),
		ecs.GetComponentType(c7{}),
	)
	entity := world.NewEntity()
	ecs.SaveComponent(world.Components(), entity, c1{})
	ecs.SaveComponent(world.Components(), entity, c2{})
	ecs.SaveComponent(world.Components(), entity, c3{})
	ecs.SaveComponent(world.Components(), entity, c4{})
	ecs.SaveComponent(world.Components(), entity, c5{})
	ecs.SaveComponent(world.Components(), entity, c6{})
	ecs.SaveComponent(world.Components(), entity, c7{})

	array := ecs.GetComponentArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.SaveComponent(entity, Component{})
		array.SaveComponent(entity, Component{})
		array.SaveComponent(entity, Component{})
		array.SaveComponent(entity, Component{})
	}
}
