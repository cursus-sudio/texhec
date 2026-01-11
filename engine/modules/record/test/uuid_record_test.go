package test

import (
	"engine/modules/uuid"
	"maps"
	"slices"
	"testing"
)

func TestUUIDForwardRecording(t *testing.T) {
	s := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 8}
	uuidComponent := uuid.New(s.UUID.NewUUID())

	entity := s.World.NewEntity()
	s.UUID.UUID().Set(entity, uuidComponent)
	s.ComponentArray.Set(entity, initialState)

	recordingID := s.Record.UUID().StartRecording(s.Config)

	s.ComponentArray.Set(entity, middleState)
	s.Record.UUID().Stop(s.Record.UUID().StartRecording(s.Config))
	s.ComponentArray.Set(entity, finalState)

	recording, ok := s.Record.UUID().Stop(recordingID)
	if !ok {
		t.Error("expected recording to exist")
		return
	}

	expected := map[uuid.UUID][]any{
		uuidComponent.ID: {finalState},
	}
	if !maps.EqualFunc(expected, recording.Entities, func(v1, v2 []any) bool {
		return slices.Equal(v1, v2)
	}) {
		t.Errorf("expected recording %v but got %v", expected, recording.Entities)
		return
	}
}

func TestUUIDBackwardsRecording(t *testing.T) {
	s := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 9}
	uuidComponent := uuid.New(s.UUID.NewUUID())

	entity := s.World.NewEntity()
	s.UUID.UUID().Set(entity, uuidComponent)
	s.ComponentArray.Set(entity, initialState)

	recordingID := s.Record.UUID().StartBackwardsRecording(s.Config)

	s.ComponentArray.Set(entity, middleState)
	s.Record.UUID().Stop(s.Record.UUID().StartRecording(s.Config))
	s.ComponentArray.Set(entity, finalState)

	recording, ok := s.Record.UUID().Stop(recordingID)
	if !ok {
		t.Error("expected recording to exist")
		return
	}

	expected := map[uuid.UUID][]any{
		uuidComponent.ID: {initialState},
	}
	if !maps.EqualFunc(expected, recording.Entities, func(v1, v2 []any) bool {
		return slices.Equal(v1, v2)
	}) {
		t.Errorf("expected recording %v but got %v", expected, recording.Entities)
		return
	}

	s.World.RemoveEntity(entity)
	s.Record.UUID().Apply(s.Config, recording)

	if c, ok := s.ComponentArray.Get(s.ComponentArray.GetEntities()[0]); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}

func TestUUIDGetState(t *testing.T) {
	s := NewSetup()
	initialState := Component{Counter: 6}

	entity := s.World.NewEntity()
	s.ComponentArray.Set(entity, initialState)

	recording := s.Record.UUID().GetState(s.Config)

	uuidComponent, ok := s.UUID.UUID().Get(entity)
	if !ok {
		t.Errorf("expected entity to get uuid component when recorded by uuid recorder")
		return
	}
	expected := map[uuid.UUID][]any{uuidComponent.ID: {initialState}}
	if !maps.EqualFunc(expected, recording.Entities, func(v1, v2 []any) bool {
		return slices.Equal(v1, v2)
	}) {
		t.Errorf("expected recording %v but got %v", expected, recording.Entities)
		return
	}
	s.ComponentArray.Remove(entity)
	s.Record.UUID().Apply(s.Config, recording)
	if ei := s.ComponentArray.GetEntities(); len(ei) != 1 {
		t.Errorf("unexpected entities on apply. expected one entity got %v", ei)
		return
	}
	// if c, ok := world.UUID().Component().Get(world.ComponentArray.GetEntities()[0]); !ok || c.ID != recording.EntitiesUUIDs.GetValues()[0] {
	// 	t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
	// 	return
	// }
	if c, ok := s.ComponentArray.Get(s.ComponentArray.GetEntities()[0]); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}
