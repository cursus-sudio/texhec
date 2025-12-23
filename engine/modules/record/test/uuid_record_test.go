package test

import (
	"engine/modules/uuid"
	"maps"
	"slices"
	"testing"
)

func TestUUIDForwardRecording(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 8}
	uuidComponent := uuid.New(world.UUID().NewUUID())

	entity := world.NewEntity()
	world.UUID().Component().Set(entity, uuidComponent)
	world.ComponentArray.Set(entity, initialState)

	recordingID := world.Record().UUID().StartRecording(world.Config)

	world.ComponentArray.Set(entity, middleState)
	world.Record().UUID().Stop(world.Record().UUID().StartRecording(world.Config))
	world.ComponentArray.Set(entity, finalState)

	recording, ok := world.Record().UUID().Stop(recordingID)
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
	world := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 9}
	uuidComponent := uuid.New(world.UUID().NewUUID())

	entity := world.NewEntity()
	world.UUID().Component().Set(entity, uuidComponent)
	world.ComponentArray.Set(entity, initialState)

	recordingID := world.Record().UUID().StartBackwardsRecording(world.Config)

	world.ComponentArray.Set(entity, middleState)
	world.Record().UUID().Stop(world.Record().UUID().StartRecording(world.Config))
	world.ComponentArray.Set(entity, finalState)

	recording, ok := world.Record().UUID().Stop(recordingID)
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

	world.RemoveEntity(entity)
	world.Record().UUID().Apply(world.Config, recording)

	if c, ok := world.ComponentArray.Get(world.ComponentArray.GetEntities()[0]); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}

func TestUUIDGetState(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, initialState)

	recording := world.Record().UUID().GetState(world.Config)

	uuidComponent, ok := world.UUID().Component().Get(entity)
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
	world.ComponentArray.Remove(entity)
	world.Record().UUID().Apply(world.Config, recording)
	if ei := world.ComponentArray.GetEntities(); len(ei) != 1 {
		t.Errorf("unexpected entities on apply. expected one entity got %v", ei)
		return
	}
	// if c, ok := world.UUID().Component().Get(world.ComponentArray.GetEntities()[0]); !ok || c.ID != recording.EntitiesUUIDs.GetValues()[0] {
	// 	t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
	// 	return
	// }
	if c, ok := world.ComponentArray.Get(world.ComponentArray.GetEntities()[0]); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}
