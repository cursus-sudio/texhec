package test

import (
	"engine/services/ecs"
	"testing"
)

func BenchmarkEntityRecording(b *testing.B) {
	s := NewSetup()

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, Component{Counter: 6})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recordingID := s.Record.Entity().StartRecording(s.Config)
		s.Record.Entity().Stop(recordingID)
	}
}

func BenchmarkCreateNEntitiesEntityRecording(b *testing.B) {
	s := NewSetup()

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, Component{Counter: 6})

	recordingID := s.Record.Entity().StartRecording(s.Config)
	for i := 0; i < b.N; i++ {
		s.ComponentArray.Set(ecs.EntityID(i), Component{Counter: i})
	}
	b.ResetTimer()
	s.Record.Entity().Stop(recordingID)
}

func BenchmarkEntityApply1Entities(b *testing.B) {
	s := NewSetup()

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, Component{Counter: 6})

	recording := s.Record.Entity().GetState(s.Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Record.Entity().Apply(s.Config, recording)
	}
}
func BenchmarkEntityApply10Entities(b *testing.B) {
	s := NewSetup()

	for i := 0; i < 10; i++ {
		entity := s.World.NewEntity()
		s.ComponentArray.Set(entity, Component{Counter: 6})
	}

	recording := s.Record.Entity().GetState(s.Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Record.Entity().Apply(s.Config, recording)
	}
}
