package test

import "testing"

func BenchmarkSpatialIndexingGetEmpty(b *testing.B) {
	setup := NewSetup()
	component := Component{Index: 69}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Tool.Get(component.Index)
	}
}

func BenchmarkSpatialIndexingGet(b *testing.B) {
	setup := NewSetup()
	component := Component{Index: 69}
	entity := setup.W.NewEntity()
	setup.Array.Set(entity, component)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Tool.Get(component.Index)
	}
}

func BenchmarkSpatialIndexingSave(b *testing.B) {
	setup := NewSetup()
	component := Component{Index: 69}
	entity := setup.W.NewEntity()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Array.Remove(entity)
		setup.Array.Set(entity, component)
	}
	setup.Tool.Get(0)
}
