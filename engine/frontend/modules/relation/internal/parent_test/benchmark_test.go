package parent_test

import "testing"

func BenchmarkGetEmpty(b *testing.B) {
	setup := NewSetup()
	tool := setup.Tool()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tool.GetChildren(0)
	}
}

func BenchmarkGet(b *testing.B) {
	setup := NewSetup()
	parent := setup.W.NewEntity()
	component := Component{Parent: parent}
	child := setup.W.NewEntity()
	if err := setup.Array.SaveComponent(child, component); err != nil {
		b.Error(err)
		return
	}
	tool := setup.Tool()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tool.GetChildren(parent)
	}
}

func BenchmarkSave(b *testing.B) {
	setup := NewSetup()
	setup.Tool()
	parent := setup.W.NewEntity()
	component := Component{Parent: parent}
	child := setup.W.NewEntity()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Array.RemoveComponent(child)
		setup.Array.SaveComponent(child, component)
	}
}

func BenchmarkSaveWithoutTool(b *testing.B) {
	setup := NewSetup()
	parent := setup.W.NewEntity()
	component := Component{Parent: parent}
	child := setup.W.NewEntity()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		setup.Array.RemoveComponent(child)
		setup.Array.SaveComponent(child, component)
	}
}
