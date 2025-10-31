package ecs_test

import (
	"shared/services/ecs"
	"testing"
)

func addQueries(world ecs.World, others ...ecs.ComponentType) {
	for _, other := range others {
		query := world.Query().Require(
			ecs.GetComponentType(Component{}),
			other,
		).Build()
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
		world.Query().Require(componentType).Build()
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

	array := ecs.GetComponentsArray[Component](world.Components())
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

	array := ecs.GetComponentsArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.SaveComponent(entity, Component{})
	}
}

func add7LiveQueriesAndAddComponenets(world ecs.World, entity ecs.EntityID) {
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
	ecs.SaveComponent(world.Components(), entity, c1{})
	ecs.SaveComponent(world.Components(), entity, c2{})
	ecs.SaveComponent(world.Components(), entity, c3{})
	ecs.SaveComponent(world.Components(), entity, c4{})
	ecs.SaveComponent(world.Components(), entity, c5{})
	ecs.SaveComponent(world.Components(), entity, c6{})
	ecs.SaveComponent(world.Components(), entity, c7{})
}

func BenchmarkSaveWith7LiveQueries(b *testing.B) {
	world := ecs.NewWorld()

	entity := world.NewEntity()
	add7LiveQueriesAndAddComponenets(world, entity)

	array := ecs.GetComponentsArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		array.SaveComponent(entity, Component{})
	}
}

func benchmarkNEntitiesSaveWith7LiveQueries(b *testing.B, entitiesCount int) {
	world := ecs.NewWorld()

	entity := world.NewEntity()
	add7LiveQueriesAndAddComponenets(world, entity)

	array := ecs.GetComponentsArray[Component](world.Components())
	array.SaveComponent(entity, Component{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < entitiesCount; j++ {
			array.SaveComponent(entity, Component{})
		}
	}
}
func benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b *testing.B, entitiesCount int) {
	world := ecs.NewWorld()

	entity := world.NewEntity()
	add7LiveQueriesAndAddComponenets(world, entity)

	array := ecs.GetComponentsArray[Component](world.Components())
	array.SaveComponent(entity, Component{})
	transaction := array.Transaction()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < entitiesCount; i++ {
			transaction.SaveComponent(entity, Component{})
		}
		transaction.Flush()
	}
}

func Benchmark4SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 4)
}
func Benchmark4SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 4)
}
func Benchmark16SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 16)
}
func Benchmark16SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 16)
}
func Benchmark256SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 256)
}
func Benchmark256SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 256)
}
func Benchmark4096SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 4096)
}
func Benchmark4096SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 4096)
}
func Benchmark16384SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 16384)
}
func Benchmark16384SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 16384)
}
func Benchmark65536SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 65536)
}
func Benchmark65536SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 65536)
}
func Benchmark262144SavesWith7LiveQueries(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueries(b, 262144)
}
func Benchmark262144SavesWith7LiveQueriesWithTransaction(b *testing.B) {
	benchmarkNEntitiesSaveWith7LiveQueriesWithTransaction(b, 262144)
}
