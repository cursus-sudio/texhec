package recordimpl

import (
	"engine/modules/record"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type Recording struct {
	Sealed datastructures.SparseSet[ecs.EntityID]
	record.Recording
}

type entityKeyedRecorder struct {
	*tool

	i     record.RecordingID
	holes datastructures.SparseSet[record.RecordingID]

	recordings          datastructures.SparseArray[record.RecordingID, Recording]
	backwardsRecordings datastructures.SparseArray[record.RecordingID, Recording]
}

func newEntityKeyedRecorder(
	t *tool,
) *entityKeyedRecorder {
	entityKeyedRecorder := &entityKeyedRecorder{
		t,

		1,
		datastructures.NewSparseSet[record.RecordingID](),

		datastructures.NewSparseArray[record.RecordingID, Recording](),
		datastructures.NewSparseArray[record.RecordingID, Recording](),
	}

	return entityKeyedRecorder
}

func (t *entityKeyedRecorder) GetState(config record.Config) record.Recording {
	recording := record.Recording{
		RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
		// Sealed:          datastructures.NewSparseSet[ecs.EntityID](),
		Arrays: make(map[string]record.ArrayRecording, len(config.RecordedComponents)),
	}
	for arrayType, arrayCtor := range config.RecordedComponents {
		array := arrayCtor(t.world)
		components := datastructures.NewSparseArray[ecs.EntityID, any]()
		recording.Arrays[arrayType.String()] = record.ArrayRecording(components)

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
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.recordings.Set(id, recording)
	return id
}
func (t *entityKeyedRecorder) StartRecording(config record.Config) record.RecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.backwardsRecordings.Set(id, recording)
	return id
}
func (t *entityKeyedRecorder) Stop(id record.RecordingID) (record.Recording, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SynchronizeState()
	if recording, ok := t.recordings.Get(id); ok {
		t.recordings.Remove(id)
		t.holes.Add(id)
		return recording.Recording, ok
	}
	if recording, ok := t.backwardsRecordings.Get(id); ok {
		t.backwardsRecordings.Remove(id)
		t.holes.Add(id)
		return recording.Recording, ok
	}
	return record.Recording{}, false
}
func (t *entityKeyedRecorder) Apply(config record.Config, recordings ...record.Recording) {
	for arrayType, arrayCtor := range config.RecordedComponents {
		t.GetArrayAndEnsureExists(arrayType, arrayCtor)
	}
	for _, recording := range recordings {
		if recording.RemovedEntities == nil ||
			recording.Arrays == nil {
			return
		}
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
		err := array.SetAny(entity, component)
		t.logger.Warn(err)
	}
}

func (t *entityKeyedRecorder) getRecordingAndID(config record.Config) (record.RecordingID, Recording) {
	var id record.RecordingID
	if holes := t.holes.GetIndices(); len(holes) != 0 {
		id = holes[0]
		t.holes.Remove(id)
	} else {
		id = t.i
		t.i++
	}

	recording := Recording{
		Sealed: datastructures.NewSparseSet[ecs.EntityID](),
		Recording: record.Recording{
			RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
			Arrays:          make(map[string]record.ArrayRecording, len(config.RecordedComponents)),
		},
	}
	for arrayType, arrayCtor := range config.RecordedComponents {
		t.GetArrayAndEnsureExists(arrayType, arrayCtor)
		recording.Arrays[arrayType.String()] = record.ArrayRecording(datastructures.NewSparseArray[ecs.EntityID, any]())
	}

	return id, recording
}
