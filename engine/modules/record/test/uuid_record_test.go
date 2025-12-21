package test

import (
	"engine/modules/uuid"
	"reflect"
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

	if len(recording.EntitiesUUIDs.GetIndices()) != 1 || recording.EntitiesUUIDs.GetIndices()[0] != entity {
		t.Errorf("expected entity [%v] got %v", entity, recording.EntitiesUUIDs.GetIndices())
		return
	}
	if len(recording.EntitiesUUIDs.GetValues()) != 1 || recording.EntitiesUUIDs.GetValues()[0] != uuidComponent.ID {
		t.Errorf("expected uuid [%v] got %v", uuidComponent.ID, recording.EntitiesUUIDs.GetValues())
		return
	}

	array, ok := recording.Arrays[reflect.TypeFor[Component]()]
	if !ok {
		t.Error("expected recording to have changes")
		return
	}
	if len(array.GetIndices()) != 1 || array.GetIndices()[0] != entity {
		t.Errorf("expected array entities to be only entity [%v] not %v", entity, array.GetIndices())
		return
	}
	if len(array.GetValues()) != 1 || array.GetValues()[0] != finalState {
		t.Errorf("expected array components to be only component [%v] not %v", finalState, array.GetValues())
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

	if len(recording.EntitiesUUIDs.GetIndices()) != 1 || recording.EntitiesUUIDs.GetIndices()[0] != entity {
		t.Errorf("expected entity [%v] got %v", entity, recording.EntitiesUUIDs.GetIndices())
		return
	}
	if len(recording.EntitiesUUIDs.GetValues()) != 1 || recording.EntitiesUUIDs.GetValues()[0] != uuidComponent.ID {
		t.Errorf("expected uuid [%v] got %v", uuidComponent.ID, recording.EntitiesUUIDs.GetValues())
		return
	}

	array, ok := recording.Arrays[reflect.TypeFor[Component]()]
	if !ok {
		t.Error("expected recording to have changes")
		return
	}
	if len(array.GetIndices()) != 1 || array.GetIndices()[0] != entity {
		t.Errorf("expected array entities to be only entity [%v] not %v", entity, array.GetIndices())
		return
	}
	if len(array.GetValues()) != 1 || array.GetValues()[0] != initialState {
		t.Errorf("expected array components to be only component [%v] not %v", initialState, array.GetValues())
		return
	}
}
