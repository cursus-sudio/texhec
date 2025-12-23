package test

import (
	"engine/services/ecs"
	"testing"
)

func BenchmarkEntityRecording(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	for i := 0; i < b.N; i++ {
		recordingID := world.Record().Entity().StartRecording(world.Config)
		world.Record().Entity().Stop(recordingID)
	}
}

func BenchmarkSetInEntityRecording(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	recordingID := world.Record().Entity().StartRecording(world.Config)
	for i := 0; i < b.N; i++ {
		world.ComponentArray.Set(entity, Component{Counter: i})
	}
	world.Record().Entity().Stop(recordingID)
}

func BenchmarkCreateNEntities(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	for i := 0; i < b.N; i++ {
		world.ComponentArray.Set(ecs.EntityID(i), Component{Counter: i})
	}
}
func BenchmarkCreateNEntitiesEntityRecording(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	recordingID := world.Record().Entity().StartRecording(world.Config)
	for i := 0; i < b.N; i++ {
		world.ComponentArray.Set(ecs.EntityID(i), Component{Counter: i})
	}
	world.Record().Entity().Stop(recordingID)
}

func BenchmarkEntityApply1Entities(b *testing.B) {
	world := NewSetup()

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, Component{Counter: 6})

	recording := world.Record().Entity().GetState(world.Config)
	for i := 0; i < b.N; i++ {
		world.Record().Entity().Apply(world.Config, recording)
	}
}
func BenchmarkEntityApply10Entities(b *testing.B) {
	world := NewSetup()

	for i := 0; i < 10; i++ {
		entity := world.NewEntity()
		world.ComponentArray.Set(entity, Component{Counter: 6})
	}

	recording := world.Record().Entity().GetState(world.Config)
	for i := 0; i < b.N; i++ {
		world.Record().Entity().Apply(world.Config, recording)
	}
}
