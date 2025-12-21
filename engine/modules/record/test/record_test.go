package test

import (
	"reflect"
	"testing"
)

func TestEntityForwardRecording(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 8}

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

func TestEntityBackwardsRecording(t *testing.T) {
	world := NewSetup()
	initialState := Component{Counter: 6}
	middleState := Component{Counter: 7}
	finalState := Component{Counter: 9}

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

	world.Record().Entity().Apply(recording)

	if ei := world.ComponentArray.GetEntities(); len(ei) != 1 || ei[0] != entity {
		t.Errorf("unexpected entities on apply. expected [%v] got %v", entity, ei)
		return
	}
	if c, ok := world.ComponentArray.Get(entity); !ok || c != initialState {
		t.Errorf("unexpected component on apply. expected %v %t got %v %t", initialState, true, c, ok)
		return
	}
}
