package state

import (
	"engine/modules/netsync"
	"engine/modules/netsync/internal/config"
	"engine/modules/uuid"
	"engine/services/ecs"
	"engine/services/logger"
)

type Tool interface {
	GetState() State
	ApplyState(State)

	StartRecording()
	FinishRecording() *State
	RecordEntitiesChange()
}

//

type toolState struct {
	recordedChanges *State
	dirtySet        ecs.DirtySet

	logger logger.Logger

	world  netsync.World
	arrays []ecs.AnyComponentArray
}

type tool struct {
	config.Config
	*toolState
}

func NewToolFactory(
	config config.Config,
	logger logger.Logger,
) ecs.ToolFactory[netsync.World, Tool] {
	// each factory client can get unique instance so mutex isn't necessary
	return ecs.NewToolFactory(func(world netsync.World) Tool {
		arrayCtors := config.ArraysOfComponents
		dirtySet := ecs.NewDirtySet()
		arrays := make([]ecs.AnyComponentArray, len(arrayCtors))
		for i, ctor := range arrayCtors {
			array := ctor(world)
			array.AddDirtySet(dirtySet)
			arrays[i] = array
		}

		t := tool{
			config,
			&toolState{
				nil,
				dirtySet,

				logger,

				world,
				arrays,
			},
		}
		return t
	})
}

func (t tool) GetState() State {
	state := State{
		Entities: make(map[uuid.UUID]EntitySnapshot),
	}
	for _, entity := range t.world.UUID().Component().GetEntities() {
		t.captureEntity(state, entity)
	}
	return state
}

func (t tool) ApplyState(changes State) {
	for id, snapshot := range changes.Entities {
		entity, ok := t.world.UUID().Entity(id)
		if snapshot.Components == nil {
			t.world.RemoveEntity(entity)
			continue
		}
		if !ok {
			entity = t.world.NewEntity()
			t.world.UUID().Component().Set(entity, uuid.New(id))
		}
		for i, array := range t.arrays {
			if snapshot.Components[i] != nil {
				err := array.SetAny(entity, snapshot.Components[i])
				t.logger.Warn(err)
			} else {
				array.Remove(entity)
			}
		}
	}
}

func (t tool) StartRecording() {
	t.recordedChanges = &State{
		Entities: map[uuid.UUID]EntitySnapshot{},
	}
}

func (t tool) FinishRecording() *State {
	changes := t.recordedChanges
	t.recordedChanges = nil
	return changes
}

func (t tool) RecordEntitiesChange() {
	recording := t.recordedChanges
	if recording == nil {
		return
	}
	for _, entity := range t.dirtySet.Get() {
		t.captureEntity(*recording, entity)
	}
}

// private methods

func (t tool) captureEntity(state State, entity ecs.EntityID) {
	unique, ok := t.world.UUID().Component().Get(entity)
	if !ok {
		return
	}

	if _, ok := state.Entities[unique.ID]; ok {
		return
	}

	snapshot := EntitySnapshot{
		Components: make([]ComponentState, len(t.Components)),
	}

	for i, array := range t.arrays {
		component, ok := array.GetAny(entity)
		if ok {
			snapshot.Components[i] = component
		}
	}

	state.Entities[unique.ID] = snapshot
}
