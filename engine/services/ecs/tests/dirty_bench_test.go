package ecs

import (
	"engine/services/ecs"
	"testing"
)

func TestDirtyFlags(t *testing.T) {
	var entity ecs.EntityID = 3
	var entity2 ecs.EntityID = 4

	set := ecs.NewDirtySet()

	if dirty := set.Get(); len(dirty) != 0 {
		t.Errorf("empty dirty flags have not expected dirty entities (%v)", dirty)
	}

	set.Dirty(entity)

	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected in dirty flags dirty entity but got %v", dirty)
	}

	if dirty := set.Get(); len(dirty) != 0 {
		t.Errorf("empty dirty flags have not expected dirty entities (%v)", dirty)
	}

	set.Dirty(entity)
	set.Dirty(entity2)

	if dirty := set.Get(); len(dirty) != 2 ||
		min(dirty[0], dirty[1]) != entity ||
		max(dirty[0], dirty[1]) != entity2 {
		t.Errorf("expected in dirty flags dirty entities but got %v", dirty)
	}
}

func TestDirtyFlagsInWorld(t *testing.T) {
	world := ecs.NewWorld()
	set := ecs.NewDirtySet()

	entity := world.NewEntity()
	type Component struct{}
	component := Component{}

	components := ecs.GetComponentsArray[Component](world)
	components.SaveComponent(entity, component)

	if dirty := set.Get(); len(dirty) != 0 {
		t.Errorf("empty dirty flags have not expected dirty entities (%v)", dirty)
		return
	}

	components.AddDirtySet(set)
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected in dirty flags dirty entity but go %v", dirty)
		return
	}

	components.RemoveComponent(entity)
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected in dirty flags dirty entity but go %v", dirty)
		return
	}

	components.SaveComponent(entity, component)
	if dirty := set.Get(); len(dirty) != 1 || dirty[0] != entity {
		t.Errorf("expected in dirty flags dirty entity but go %v", dirty)
		return
	}
}

//

func benchmarkNEntitiesSaveWith7Systems(b *testing.B, entitiesCount int) {
	type Component1 struct{}
	type Component2 struct{}
	type Component3 struct{}
	type Component4 struct{}
	type Component5 struct{}
	type Component6 struct{}
	type Component7 struct{}

	world := ecs.NewWorld()
	set := ecs.NewDirtySet()

	arr1 := ecs.GetComponentsArray[Component1](world)
	arr2 := ecs.GetComponentsArray[Component2](world)
	arr3 := ecs.GetComponentsArray[Component3](world)
	arr4 := ecs.GetComponentsArray[Component4](world)
	arr5 := ecs.GetComponentsArray[Component5](world)
	arr6 := ecs.GetComponentsArray[Component6](world)
	arr7 := ecs.GetComponentsArray[Component7](world)

	arr1.AddDirtySet(set)
	arr2.AddDirtySet(set)
	arr3.AddDirtySet(set)
	arr4.AddDirtySet(set)
	arr5.AddDirtySet(set)
	arr6.AddDirtySet(set)
	arr7.AddDirtySet(set)

	entity := world.NewEntity()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < entitiesCount; j++ {
			arr1.SaveComponent(entity, Component1{})
		}
	}
}

func Benchmark4SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 4)
}
func Benchmark16SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 16)
}
func Benchmark256SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 256)
}
func Benchmark4096SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 4096)
}
func Benchmark16384SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 16384)
}
func Benchmark65536SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 65536)
}
func Benchmark262144SavesWith7Systems(b *testing.B) {
	benchmarkNEntitiesSaveWith7Systems(b, 262144)
}
