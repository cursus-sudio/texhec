package recordimpl

import (
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"fmt"
	"reflect"
)

type uuidKeyedRecorder struct {
	*tool

	i     record.UUIDRecordingID
	holes datastructures.SparseSet[record.UUIDRecordingID]

	uuidRecordings          datastructures.SparseArray[record.UUIDRecordingID, record.UUIDRecording]
	backwardsUUIDRecordings datastructures.SparseArray[record.UUIDRecordingID, record.UUIDRecording]
}

func newUUIDKeyedRecorder(
	t *tool,
) *uuidKeyedRecorder {
	uuidKeyedRecorder := &uuidKeyedRecorder{
		t,

		0,
		datastructures.NewSparseSet[record.UUIDRecordingID](),

		datastructures.NewSparseArray[record.UUIDRecordingID, record.UUIDRecording](),
		datastructures.NewSparseArray[record.UUIDRecordingID, record.UUIDRecording](),
	}

	return uuidKeyedRecorder
}

func (t *uuidKeyedRecorder) GetState(config record.Config) record.UUIDRecording {
	recording := record.UUIDRecording{
		UUIDEntities:  make(map[uuid.UUID]ecs.EntityID),
		EntitiesUUIDs: datastructures.NewSparseArray[ecs.EntityID, uuid.UUID](),
		Recording: record.Recording{
			RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
			Sealed:          datastructures.NewSparseSet[ecs.EntityID](),
			Arrays:          make(map[reflect.Type]record.ArrayRecording, len(config.RecordedComponents)),
		},
	}
	entitiesUUIDs := datastructures.NewSparseArray[ecs.EntityID, uuid.UUID]()
	for arrayType, arrayCtor := range config.RecordedComponents {
		array := arrayCtor(t.world)
		components := datastructures.NewSparseArray[ecs.EntityID, any]()
		recording.Arrays[arrayType] = components

		//

		for _, entity := range array.GetEntities() {
			if _, ok := entitiesUUIDs.Get(entity); !ok {
				uuid, ok := t.world.UUID().Component().Get(entity)
				if !ok {
					uuid.ID = t.world.UUID().NewUUID()
					t.world.UUID().Component().Set(entity, uuid)
				}
				entitiesUUIDs.Set(entity, uuid.ID)

			}
			component, ok := array.GetAny(entity)
			if !ok {
				continue
			}
			components.Set(entity, component)
		}
	}
	return recording
}

func (t *uuidKeyedRecorder) StartBackwardsRecording(config record.Config) record.UUIDRecordingID {
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.uuidRecordings.Set(id, recording)
	return id
}
func (t *uuidKeyedRecorder) StartRecording(config record.Config) record.UUIDRecordingID {
	t.SynchronizeState()
	id, recording := t.getRecordingAndID(config)
	t.backwardsUUIDRecordings.Set(id, recording)
	return id
}
func (t *uuidKeyedRecorder) Stop(id record.UUIDRecordingID) (record.UUIDRecording, bool) {
	t.SynchronizeState()
	if recording, ok := t.uuidRecordings.Get(id); ok {
		t.uuidRecordings.Remove(id)
		t.holes.Add(id)
		return recording, ok
	}
	if recording, ok := t.backwardsUUIDRecordings.Get(id); ok {
		t.backwardsUUIDRecordings.Remove(id)
		t.holes.Add(id)
		return recording, ok
	}
	return record.UUIDRecording{}, false
}
func (t *uuidKeyedRecorder) Apply(recordings ...record.UUIDRecording) {
	errs := []error{}
	for _, recording := range recordings {
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
func (t *uuidKeyedRecorder) getRecordingAndID(config record.Config) (record.UUIDRecordingID, record.UUIDRecording) {
	var id record.UUIDRecordingID
	if holes := t.holes.GetIndices(); len(holes) != 0 {
		id = holes[0]
		t.holes.Remove(id)
	} else {
		id = t.i
		t.i++
	}

	recording := record.UUIDRecording{
		UUIDEntities:  make(map[uuid.UUID]ecs.EntityID),
		EntitiesUUIDs: datastructures.NewSparseArray[ecs.EntityID, uuid.UUID](),
		Recording: record.Recording{
			RemovedEntities: datastructures.NewSparseSet[ecs.EntityID](),
			Sealed:          datastructures.NewSparseSet[ecs.EntityID](),
			Arrays:          make(map[reflect.Type]record.ArrayRecording, len(config.RecordedComponents)),
		},
	}

	for arrayType, arrayCtor := range config.RecordedComponents {
		t.tool.GetArrayAndEnsureExists(arrayType, arrayCtor)
		recording.Arrays[arrayType] = datastructures.NewSparseArray[ecs.EntityID, any]()
	}
	return id, recording
}
