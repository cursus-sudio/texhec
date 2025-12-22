package test

import (
	"engine/modules/uuid"
	"engine/services/ecs"
	"maps"
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

	if expected := map[uuid.UUID]ecs.EntityID{uuidComponent.ID: entity}; !maps.Equal(recording.UUIDEntities, expected) {
		t.Errorf("expected [%v] got [%v]", expected, recording.UUIDEntities)
		return
	}

	array, ok := recording.Arrays[reflect.TypeFor[Component]().String()]
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

	if expected := map[uuid.UUID]ecs.EntityID{uuidComponent.ID: entity}; !maps.Equal(recording.UUIDEntities, expected) {
		t.Errorf("expected [%v] got [%v]", expected, recording.UUIDEntities)
		return
	}

	array, ok := recording.Arrays[reflect.TypeFor[Component]().String()]
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

	world.RemoveEntity(entity)
	world.Record().UUID().Apply(world.Config, recording)

	if ei := world.ComponentArray.GetEntities(); len(ei) != 1 {
		t.Errorf("unexpected entities on apply. expected one entity got %v", ei)
		return
	}
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

	array, ok := recording.Arrays[reflect.TypeFor[Component]().String()]
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

	if uuidComponent, ok := world.UUID().Component().Get(entity); !ok {
		t.Errorf("entity should have and uuid component")
		return
	} else if _, ok := recording.UUIDEntities[uuidComponent.ID]; !ok {
		t.Errorf("entity should have and uuid in recording")
		return
	}
	array.Remove(entity)
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
