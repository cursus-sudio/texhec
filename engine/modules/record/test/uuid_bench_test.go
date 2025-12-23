package test

import (
	"engine/services/ecs"
	"testing"
)

func BenchmarkUUIDRecording(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recordingID := world.Record().UUID().StartRecording(world.Config)
		world.Record().UUID().Stop(recordingID)
	}
}

func BenchmarkCreateNEntitiesUUIDRecording(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	recordingID := world.Record().UUID().StartRecording(world.Config)
	for i := 0; i < b.N; i++ {
		world.ComponentArray.Set(ecs.EntityID(i), Component{Counter: i})
	}
	b.ResetTimer()
	world.Record().UUID().Stop(recordingID)
}

func BenchmarkUUIDApply(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	recording := world.Record().UUID().GetState(world.Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Record().UUID().Apply(world.Config, recording)
	}
}
func BenchmarkUUIDApply10Entities(b *testing.B) {
	world := NewSetup()

	for i := 0; i < 10; i++ {
		entity := world.NewEntity()
		world.ComponentArray.Set(entity, Component{Counter: 6})
	}

	recording := world.Record().UUID().GetState(world.Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Record().UUID().Apply(world.Config, recording)
	}
}
