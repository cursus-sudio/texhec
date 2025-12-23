package recordimpl

import (
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type UUIDForwardRecording struct {
	Config   record.Config
	DirtySet ecs.DirtySet
}
type UUIDBackwardRecording struct {
	Config          record.Config
	WorldCopyArrays []ecs.AnyComponentArray
	Entities        map[uuid.UUID][]any
}

type uuidKeyedRecorder struct {
	*tool

	i     record.UUIDRecordingID
	holes datastructures.SparseSet[record.UUIDRecordingID]

	forwardRecordings   datastructures.SparseArray[record.UUIDRecordingID, *UUIDForwardRecording]
	backwardsRecordings datastructures.SparseArray[record.UUIDRecordingID, *UUIDBackwardRecording]
}

func newUUIDKeyedRecorder(
	t *tool,
) *uuidKeyedRecorder {
	uuidKeyedRecorder := &uuidKeyedRecorder{
		t,

		1,
		datastructures.NewSparseSet[record.UUIDRecordingID](),

		datastructures.NewSparseArray[record.UUIDRecordingID, *UUIDForwardRecording](),
		datastructures.NewSparseArray[record.UUIDRecordingID, *UUIDBackwardRecording](),
	}

	return uuidKeyedRecorder
}

func (t *uuidKeyedRecorder) GetState(config record.Config) record.UUIDRecording {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		for _, entity := range array.GetEntities() {
			entities.Add(entity)
		}
	}
	return t.getStateFor(config, entities.GetIndices())
}

func (t *uuidKeyedRecorder) StartBackwardsRecording(config record.Config) record.UUIDRecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	defer t.SyncBackwardsRecordingState()

	id := t.getID()
	recording := &UUIDBackwardRecording{
		Config:          config,
		WorldCopyArrays: make([]ecs.AnyComponentArray, 0, len(*config.ComponentsOrder)),
		Entities:        map[uuid.UUID][]any{},
	}
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		array.uuidDependencies.Add(recording)
		worldCopyArray := t.GetWorldCopyArray(arrayType, config)
		recording.WorldCopyArrays = append(recording.WorldCopyArrays, worldCopyArray)
	}
	t.backwardsRecordings.Set(id, recording)

	return id
}
func (t *uuidKeyedRecorder) StartRecording(config record.Config) record.UUIDRecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SyncBackwardsRecordingState()

	id := t.getID()
	recording := &UUIDForwardRecording{
		Config:   config,
		DirtySet: ecs.NewDirtySet(),
	}
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		array.AddDirtySet(recording.DirtySet)
	}
	t.forwardRecordings.Set(id, recording)

	return id
}
func (t *uuidKeyedRecorder) Stop(id record.UUIDRecordingID) (record.UUIDRecording, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if recording, ok := t.forwardRecordings.Get(id); ok {
		entities := recording.DirtySet.Get()
		recording.DirtySet.Release()
		t.forwardRecordings.Remove(id)
		t.holes.Add(id)
		return t.getStateFor(recording.Config, entities), true
	}
	if recording, ok := t.backwardsRecordings.Get(id); ok {
		t.SyncBackwardsRecordingState()
		for _, arrayType := range *recording.Config.ComponentsOrder {
			array := t.GetWorldArray(arrayType, recording.Config)
			array.uuidDependencies.RemoveElements(recording)
		}
		t.backwardsRecordings.Remove(id)
		t.holes.Add(id)
		return record.UUIDRecording{Entities: recording.Entities}, true
	}
	return record.UUIDRecording{}, false
}
func (t *uuidKeyedRecorder) Apply(config record.Config, recordings ...record.UUIDRecording) {
	arrays := make([]ecs.AnyComponentArray, 0, len(*config.ComponentsOrder))
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		arrays = append(arrays, array)
	}

	for _, recording := range recordings {
		for uuidValue, components := range recording.Entities {
			entity, ok := t.world.UUID().Entity(uuidValue)
			if !ok && components != nil {
				entity = t.world.NewEntity()
				t.world.UUID().Component().Set(entity, uuid.New(uuidValue))
			}
			if components == nil {
				t.world.RemoveEntity(entity)
				continue
			}

			for i, component := range components {
				array := arrays[i]
				if component == nil {
					array.Remove(entity)
					continue
				}
				_ = array.SetAny(entity, component)
			}
		}
	}
}

func (t *uuidKeyedRecorder) getStateFor(config record.Config, entities []ecs.EntityID) record.UUIDRecording {
	arrays := make([]ecs.AnyComponentArray, 0, len(*config.ComponentsOrder))
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		arrays = append(arrays, array)
	}

	recording := record.UUIDRecording{
		Entities: make(map[uuid.UUID][]any, len(entities)),
	}

	for _, entity := range entities {
		uuidComponent, ok := t.world.UUID().Component().Get(entity)
		if !ok {
			uuidComponent.ID = t.world.UUID().NewUUID()
			t.world.UUID().Component().Set(entity, uuidComponent)
		}
		components := make([]any, 0, len(arrays))
		for _, array := range arrays {
			v, ok := array.GetAny(entity)
			if !ok {
				v = nil
			}
			components = append(components, v)
		}
		recording.Entities[uuidComponent.ID] = components
	}

	return recording
}

func (t *uuidKeyedRecorder) getID() record.UUIDRecordingID {
	if holes := t.holes.GetIndices(); len(holes) != 0 {
		hole := holes[0]
		t.holes.Remove(hole)
		return hole
	}
	i := t.i
	t.i++
	return i
}
