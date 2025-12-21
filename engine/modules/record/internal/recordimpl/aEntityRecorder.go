package recordimpl

import (
	"engine/modules/record"
	"engine/services/datastructures"
	"engine/services/ecs"
	"reflect"
)

type entityKeyedRecorder struct {
	*tool

	i     record.RecordingID
	holes datastructures.SparseSet[record.RecordingID]

	recordings          datastructures.SparseArray[record.RecordingID, record.Recording]
	backwardsRecordings datastructures.SparseArray[record.RecordingID, record.Recording]
}

func newEntityKeyedRecorder(
	t *tool,
) *entityKeyedRecorder {
	entityKeyedRecorder := &entityKeyedRecorder{
		t,

		0,
		datastructures.NewSparseSet[record.RecordingID](),

		datastructures.NewSparseArray[record.RecordingID, record.Recording](),
		datastructures.NewSparseArray[record.RecordingID, record.Recording](),
	}

	return entityKeyedRecorder
}

func (t *entityKeyedRecorder) GetState(config record.Config) record.Recording {
	recording := record.Recording{
		RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
		Sealed:          datastructures.NewSparseSet[ecs.EntityID](),
		Arrays:          make(map[reflect.Type]record.ArrayRecording, len(config.RecordedComponents)),
	}
	for arrayType, arrayCtor := range config.RecordedComponents {
		array := arrayCtor(t.world)
		components := datastructures.NewSparseArray[ecs.EntityID, any]()
		recording.Arrays[arrayType] = components

		//

		for _, entity := range array.GetEntities() {
			component, ok := array.GetAny(entity)
			if !ok {
				continue
			}
			components.Set(entity, component)
		}
	}
	return recording
}
func (t *entityKeyedRecorder) StartBackwardsRecording(config record.Config) record.RecordingID {
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.recordings.Set(id, recording)
	return id
}
func (t *entityKeyedRecorder) StartRecording(config record.Config) record.RecordingID {
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.backwardsRecordings.Set(id, recording)
	return id
}
func (t *entityKeyedRecorder) Stop(id record.RecordingID) (record.Recording, bool) {
	t.SynchronizeState()
	if recording, ok := t.recordings.Get(id); ok {
		t.recordings.Remove(id)
		t.holes.Add(id)
		return recording, ok
	}
	if recording, ok := t.backwardsRecordings.Get(id); ok {
		t.backwardsRecordings.Remove(id)
		t.holes.Add(id)
		return recording, ok
	}
	return record.Recording{}, false
}
func (t *entityKeyedRecorder) Apply(recordings ...record.Recording) {
	for _, recording := range recordings {
		for _, entity := range recording.RemovedEntities.GetIndices() {
			t.world.RemoveEntity(entity)
		}
		for arrayType, arrayData := range recording.Arrays {
			t.applyArray(t.worldArrays[arrayType], arrayData)
		}
	}
}

//

func (t *entityKeyedRecorder) applyArray(array ecs.AnyComponentArray, arrayData record.ArrayRecording) {
	for _, entity := range arrayData.GetIndices() {
		component, _ := arrayData.Get(entity)
		if component == nil {
			array.Remove(entity)
			continue
		}
		array.SetAny(entity, component)
	}
}

func (t *entityKeyedRecorder) getRecordingAndID(config record.Config) (record.RecordingID, record.Recording) {
	var id record.RecordingID
	if holes := t.holes.GetIndices(); len(holes) != 0 {
		id = holes[0]
		t.holes.Remove(id)
	} else {
		id = t.i
		t.i++
	}

	recording := record.Recording{
		RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
		Sealed:          datastructures.NewSparseSet[ecs.EntityID](),
		Arrays:          make(map[reflect.Type]record.ArrayRecording, len(config.RecordedComponents)),
	}
	for arrayType, arrayCtor := range config.RecordedComponents {
		t.tool.GetArrayAndEnsureExists(arrayType, arrayCtor)
		recording.Arrays[arrayType] = datastructures.NewSparseArray[ecs.EntityID, any]()
	}

	return id, recording
}
