package test

import (
	"testing"
)

func TestEntityForwardRecording(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 8}

	world.ComponentArray.Set(world.NewEntity(), initialState)

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, initialState)

	recordingID := world.Record().Entity().StartRecording(world.Config)

	world.ComponentArray.Set(entity, middleState)
	world.Record().Entity().Stop(world.Record().Entity().StartRecording(world.Config))
	world.ComponentArray.Set(entity, finalState)

	recording, ok := world.Record().Entity().Stop(recordingID)
	if !ok {
		t.Error("expected recording to exist")
		return
	}

	if len(recording.Entities.GetIndices()) != 1 || recording.Entities.GetIndices()[0] != entity {
		t.Errorf("expected array entities to be only entity [%v] not %v", entity, recording.Entities.GetIndices())
		return
	}
	if len(recording.Entities.GetValues()) != 1 ||
		len(recording.Entities.GetValues()[0]) != 1 ||
		recording.Entities.GetValues()[0][0] != finalState {
		t.Errorf("expected array components to be only component [%v] not %v", finalState, recording.Entities.GetValues())
		return
	}
}

func TestEntityBackwardsRecording(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 9}

	world.ComponentArray.Set(world.NewEntity(), initialState)

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, initialState)

	recordingID := world.Record().Entity().StartBackwardsRecording(world.Config)

	world.ComponentArray.Set(entity, middleState)
	world.Record().Entity().Stop(world.Record().Entity().StartRecording(world.Config))
	world.ComponentArray.Set(entity, finalState)

	recording, ok := world.Record().Entity().Stop(recordingID)
	if !ok {
		t.Error("expected recording to exist")
		return
	}

	if len(recording.Entities.GetIndices()) != 1 || recording.Entities.GetIndices()[0] != entity {
		t.Errorf("expected array entities to be only entity [%v] not %v", entity, recording.Entities.GetIndices())
		return
	}
	if len(recording.Entities.GetValues()) != 1 ||
		len(recording.Entities.GetValues()[0]) != 1 ||
		recording.Entities.GetValues()[0][0] != initialState {
		t.Errorf("expected array components to be only component [%v] not %v", initialState, recording.Entities.GetValues())
		return
	}

	world.Record().Entity().Apply(world.Config, recording)

	if c, ok := world.ComponentArray.Get(entity); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}

func TestEntityGetState(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}

	entity := world.NewEntity()
	world.ComponentArray.Set(entity, initialState)

	recording := world.Record().Entity().GetState(world.Config)

	if len(recording.Entities.GetIndices()) != 1 || recording.Entities.GetIndices()[0] != entity {
		t.Errorf("expected array entities to be only entity [%v] not %v", entity, recording.Entities.GetIndices())
		return
	}
	if len(recording.Entities.GetValues()) != 1 ||
		len(recording.Entities.GetValues()[0]) != 1 ||
		recording.Entities.GetValues()[0][0] != initialState {
		t.Errorf("expected array components to be only component [%v] not %v", initialState, recording.Entities.GetValues())
		return
	}

	world.ComponentArray.Remove(entity)
	world.Record().Entity().Apply(world.Config, recording)
	if ei := world.ComponentArray.GetEntities(); len(ei) != 1 {
		t.Errorf("unexpected entities on apply. expected one entity got %v", ei)
		return
	}
	if c, ok := world.ComponentArray.Get(world.ComponentArray.GetEntities()[0]); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}
