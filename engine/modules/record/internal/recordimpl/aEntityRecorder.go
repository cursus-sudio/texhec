package recordimpl

import (
	"engine/modules/record"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type FowardRecording struct {
	Config   record.Config
	DirtySet ecs.DirtySet
}

type BackwardRecording struct {
	Config          record.Config
	WorldCopyArrays []ecs.AnyComponentArray
	Entities        datastructures.SparseArray[ecs.EntityID, []any]
}

type entityKeyedRecorder struct {
	*service

	i     record.RecordingID
	holes datastructures.SparseSet[record.RecordingID]

	forwardRecordings   datastructures.SparseArray[record.RecordingID, *FowardRecording]
	backwardsRecordings datastructures.SparseArray[record.RecordingID, *BackwardRecording]
}

func newEntityKeyedRecorder(
	t *service,
) *entityKeyedRecorder {
	entityKeyedRecorder := &entityKeyedRecorder{
		t,

		1,
		datastructures.NewSparseSet[record.RecordingID](),

		datastructures.NewSparseArray[record.RecordingID, *FowardRecording](),
		datastructures.NewSparseArray[record.RecordingID, *BackwardRecording](),
	}

	return entityKeyedRecorder
}

func (t *entityKeyedRecorder) GetState(config record.Config) record.Recording {
	entities := datastructures.NewSparseSet[ecs.EntityID]()
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		for _, entity := range array.GetEntities() {
			entities.Add(entity)
		}
	}
	return t.getStateFor(config, entities.GetIndices())
}
func (t *entityKeyedRecorder) StartBackwardsRecording(config record.Config) record.RecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SyncBackwardsRecordingState()

	id := t.getID()
	recording := &BackwardRecording{
		Config:          config,
		WorldCopyArrays: make([]ecs.AnyComponentArray, 0, len(*config.ComponentsOrder)),
		Entities:        datastructures.NewSparseArray[ecs.EntityID, []any](),
	}
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		array.dependencies.Add(recording)
		worldCopyArray := t.GetWorldCopyArray(arrayType, config)
		recording.WorldCopyArrays = append(recording.WorldCopyArrays, worldCopyArray)
	}
	t.backwardsRecordings.Set(id, recording)

	return id
}
func (t *entityKeyedRecorder) StartRecording(config record.Config) record.RecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	id := t.getID()
	recording := &FowardRecording{
		Config:   config,
		DirtySet: ecs.NewDirtySet(),
	}
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		array.AddDirtySet(recording.DirtySet)
	}
	recording.DirtySet.Clear()
	t.forwardRecordings.Set(id, recording)

	return id
}
func (t *entityKeyedRecorder) Stop(id record.RecordingID) (record.Recording, bool) {
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
			array.dependencies.RemoveElements(recording)
		}

		t.backwardsRecordings.Remove(id)
		t.holes.Add(id)

		return record.Recording{Entities: recording.Entities}, true
	}
	return record.Recording{}, false
}
func (t *entityKeyedRecorder) Apply(config record.Config, recordings ...record.Recording) {
	arrays := make([]ecs.AnyComponentArray, 0, len(*config.ComponentsOrder))
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		arrays = append(arrays, array)
	}

	for _, recording := range recordings {
		for _, entity := range recording.Entities.GetIndices() {
			components, ok := recording.Entities.Get(entity)
			if !ok {
				continue
			}
			if components == nil {
				t.World.RemoveEntity(entity)
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

func (t *entityKeyedRecorder) getStateFor(config record.Config, entities []ecs.EntityID) record.Recording {
	arrays := make([]ecs.AnyComponentArray, 0, len(*config.ComponentsOrder))
	for _, arrayType := range *config.ComponentsOrder {
		array := t.GetWorldArray(arrayType, config)
		arrays = append(arrays, array)
	}

	recording := record.Recording{
		Entities: datastructures.NewSparseArray[ecs.EntityID, []any](),
	}

	for _, entity := range entities {
		components := make([]any, 0, len(arrays))
		for _, array := range arrays {
			v, ok := array.GetAny(entity)
			if !ok {
				v = nil
			}
			components = append(components, v)
		}
		recording.Entities.Set(entity, components)
	}

	return recording
}

//

func (t *entityKeyedRecorder) getID() record.RecordingID {
	if holes := t.holes.GetIndices(); len(holes) != 0 {
		hole := holes[0]
		t.holes.Remove(hole)
		return hole
	}
	i := t.i
	t.i++
	return i
}
