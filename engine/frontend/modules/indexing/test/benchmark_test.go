package test

import "testing"

func BenchmarkSpatialIndexingGetEmpty(b *testing.B) {
	setup := NewSetup()
	component := Component{Index: 69}
	tool := setup.Tool()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tool.Get(component.Index)
	}
}

func BenchmarkSpatialIndexingGet(b *testing.B) {
	setup := NewSetup()
	component := Component{Index: 69}
	entity := setup.W.NewEntity()
	if err := setup.Array.SaveComponent(entity, component); err != nil {
		b.Error(err)
		return
	}
	tool := setup.Tool()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tool.Get(component.Index)
	}
}

func BenchmarkSpatialIndexingSave(b *testing.B) {
	setup := NewSetup()
	setup.Tool()
	component := Component{Index: 69}
	entity := setup.W.NewEntity()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Array.SaveComponent(entity, component)
	}
}

func BenchmarkSpatialIndexingSaveWithoutTool(b *testing.B) {
	setup := NewSetup()
	component := Component{Index: 69}
	entity := setup.W.NewEntity()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Array.SaveComponent(entity, component)
	}
}
