package test

import (
	"engine/services/ecs"
	"testing"
)

func BenchmarkUUIDRecording(b *testing.B) {
	s := NewSetup()

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, Component{Counter: 6})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recordingID := s.Record.UUID().StartRecording(s.Config)
		s.Record.UUID().Stop(recordingID)
	}
}

func BenchmarkCreateNEntitiesUUIDRecording(b *testing.B) {
	s := NewSetup()

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, Component{Counter: 6})

	recordingID := s.Record.UUID().StartRecording(s.Config)
	for i := 0; i < b.N; i++ {
		s.ComponentArray.Set(ecs.EntityID(i), Component{Counter: i})
	}
	b.ResetTimer()
	s.Record.UUID().Stop(recordingID)
}

func BenchmarkUUIDApply(b *testing.B) {
	s := NewSetup()

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, Component{Counter: 6})

	recording := s.Record.UUID().GetState(s.Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Record.UUID().Apply(s.Config, recording)
	}
}
func BenchmarkUUIDApply10Entities(b *testing.B) {
	s := NewSetup()

	for i := 0; i < 10; i++ {
		entity := s.World.NewEntity()
		s.ComponentArray.Set(entity, Component{Counter: 6})
	}

	recording := s.Record.UUID().GetState(s.Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Record.UUID().Apply(s.Config, recording)
	}
}
