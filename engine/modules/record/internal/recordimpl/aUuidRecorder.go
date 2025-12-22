package recordimpl

import (
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"fmt"
)

type UUIDRecording struct {
	Sealed        datastructures.SparseSet[ecs.EntityID]
	EntitiesUUIDs datastructures.SparseArray[ecs.EntityID, uuid.UUID]
	record.UUIDRecording
}

type uuidKeyedRecorder struct {
	*tool

	i     record.UUIDRecordingID
	holes datastructures.SparseSet[record.UUIDRecordingID]

	uuidRecordings          datastructures.SparseArray[record.UUIDRecordingID, UUIDRecording]
	backwardsUUIDRecordings datastructures.SparseArray[record.UUIDRecordingID, UUIDRecording]
}

func newUUIDKeyedRecorder(
	t *tool,
) *uuidKeyedRecorder {
	uuidKeyedRecorder := &uuidKeyedRecorder{
		t,

		1,
		datastructures.NewSparseSet[record.UUIDRecordingID](),

		datastructures.NewSparseArray[record.UUIDRecordingID, UUIDRecording](),
		datastructures.NewSparseArray[record.UUIDRecordingID, UUIDRecording](),
	}

	return uuidKeyedRecorder
}

func (t *uuidKeyedRecorder) GetState(config record.Config) record.UUIDRecording {
	recording := UUIDRecording{
		EntitiesUUIDs: datastructures.NewSparseArray[ecs.EntityID, uuid.UUID](),
		UUIDRecording: record.UUIDRecording{
			UUIDEntities: make(map[uuid.UUID]ecs.EntityID),
			Recording: record.Recording{
				RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
				Arrays:          make(map[string]record.ArrayRecording, len(config.RecordedComponents)),
			},
		},
	}
	for arrayType, arrayCtor := range config.RecordedComponents {
		array := arrayCtor(t.world)
		components := datastructures.NewSparseArray[ecs.EntityID, any]()
		recording.Arrays[arrayType.String()] = record.ArrayRecording(components)

		//

		for _, entity := range array.GetEntities() {
			if _, ok := recording.EntitiesUUIDs.Get(entity); !ok {
				uuid, ok := t.world.UUID().Component().Get(entity)
				if !ok {
					uuid.ID = t.world.UUID().NewUUID()
					t.world.UUID().Component().Set(entity, uuid)
				}
				recording.EntitiesUUIDs.Set(entity, uuid.ID)
				recording.UUIDEntities[uuid.ID] = entity
			}
			component, ok := array.GetAny(entity)
			if !ok {
				continue
			}
			components.Set(entity, component)
		}
	}
	return recording.UUIDRecording
}

func (t *uuidKeyedRecorder) StartBackwardsRecording(config record.Config) record.UUIDRecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.uuidRecordings.Set(id, recording)
	return id
}
func (t *uuidKeyedRecorder) StartRecording(config record.Config) record.UUIDRecordingID {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.backwardsUUIDRecordings.Set(id, recording)
	return id
}
func (t *uuidKeyedRecorder) Stop(id record.UUIDRecordingID) (record.UUIDRecording, bool) {
	// t.logger.Info("stopped recording")
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.SynchronizeState()
	if recording, ok := t.uuidRecordings.Get(id); ok {
		t.uuidRecordings.Remove(id)
		t.holes.Add(id)
		return recording.UUIDRecording, ok
	}
	if recording, ok := t.backwardsUUIDRecordings.Get(id); ok {
		t.backwardsUUIDRecordings.Remove(id)
		t.holes.Add(id)
		return recording.UUIDRecording, ok
	}
	return record.UUIDRecording{}, false
}
func (t *uuidKeyedRecorder) Apply(config record.Config, recordings ...record.UUIDRecording) {
	for arrayType, arrayCtor := range config.RecordedComponents {
		t.tool.GetArrayAndEnsureExists(arrayType, arrayCtor)
	}
	errs := []error{}
	for _, recording := range recordings {
		if recording.RemovedEntities == nil ||
			recording.Arrays == nil ||
			recording.UUIDEntities == nil {
			return
		}
		// entities uuids mapping
		entities := datastructures.NewSparseArray[ecs.EntityID, ecs.EntityID]()
		for uuidValue, entityPlaceholder := range recording.UUIDEntities {
			entity, ok := t.world.UUID().Entity(uuidValue)
			if !ok {
				entity = t.world.NewEntity()
				t.world.UUID().Component().Set(entity, uuid.New(uuidValue))
			}
			entities.Set(entityPlaceholder, entity)
		}

		// applying changes
		for _, entityPlaceholder := range recording.RemovedEntities.GetIndices() {
			entity, ok := entities.Get(entityPlaceholder)
			if !ok {
				continue
			}
			t.world.RemoveEntity(entity)
		}
		for arrayType, arrayData := range recording.Arrays {
			errs = append(errs, t.applyArray(
				entities,
				t.worldArrays[arrayType],
				arrayData,
			)...)
		}
	}
	if len(errs) == 0 {
		return
	}
	t.logger.Warn(fmt.Errorf(
		"error when parsing uuids and entities placeholders to entities %v",
		recordings,
	))
	for _, err := range errs {
		t.logger.Warn(err)
	}
}

func (t *uuidKeyedRecorder) applyArray(
	entities datastructures.SparseArray[ecs.EntityID, ecs.EntityID],
	array ecs.AnyComponentArray,
	arrayData record.ArrayRecording,
) []error {
	errs := []error{}
	for _, entityPlaceholder := range arrayData.GetIndices() {
		entity, ok := entities.Get(entityPlaceholder)
		if !ok {
			errs = append(errs, fmt.Errorf("entity placeholder %v lacks entity", entityPlaceholder))
			continue
		}
		component, _ := arrayData.Get(entityPlaceholder)
		if component == nil {
			array.Remove(entity)
			continue
		}
		array.SetAny(entity, component)
	}
	return errs
}
func (t *uuidKeyedRecorder) getRecordingAndID(config record.Config) (record.UUIDRecordingID, UUIDRecording) {
	var id record.UUIDRecordingID
	if holes := t.holes.GetIndices(); len(holes) != 0 {
		id = holes[0]
		t.holes.Remove(id)
	} else {
		id = t.i
		t.i++
	}

	recording := UUIDRecording{
		Sealed:        datastructures.NewSparseSet[ecs.EntityID](),
		EntitiesUUIDs: datastructures.NewSparseArray[ecs.EntityID, uuid.UUID](),
		UUIDRecording: record.UUIDRecording{
			UUIDEntities: make(map[uuid.UUID]ecs.EntityID),
			Recording: record.Recording{
				RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
				Arrays:          make(map[string]record.ArrayRecording, len(config.RecordedComponents)),
			},
		},
	}

	for arrayType, arrayCtor := range config.RecordedComponents {
		t.tool.GetArrayAndEnsureExists(arrayType, arrayCtor)
		recording.Arrays[arrayType.String()] = record.ArrayRecording(datastructures.NewSparseArray[ecs.EntityID, any]())
	}
	return id, recording
}
