package recordimpl

import (
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"reflect"
	"sync"
)

type entityArray struct {
	dirtySet ecs.DirtySet
	ecs.AnyComponentArray
}

type tool struct {
	world       record.World
	worldArrays map[string]entityArray

	worldCopy       record.World
	worldCopyArrays map[string]entityArray

	logger logger.Logger
	mutex  *sync.Mutex
	*entityKeyedRecorder
	*uuidKeyedRecorder
}

func NewToolFactory(
	uuidToolFactory uuid.ToolFactory,
	logger logger.Logger,
) record.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w record.World) record.RecordTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
			world:       w,
			worldArrays: make(map[string]entityArray),

			worldCopy:       newWorld(uuidToolFactory),
			worldCopyArrays: make(map[string]entityArray),

			logger: logger,
			mutex:  &sync.Mutex{},
		}
		t.entityKeyedRecorder = newEntityKeyedRecorder(&t)
		t.uuidKeyedRecorder = newUUIDKeyedRecorder(&t)
		w.SaveGlobal(t)

		return t
	})
}

func (t tool) Record() record.Interface {
	return t
}

//

func (t tool) Entity() record.EntityKeyedRecorder {
	return t.entityKeyedRecorder
}
func (t tool) UUID() record.UUIDKeyedRecorder {
	return t.uuidKeyedRecorder
}

//

func (t tool) GetArrayAndEnsureExists(arrayType reflect.Type, arrayCtor func(ecs.World) ecs.AnyComponentArray) {
	worldArray := entityArray{
		ecs.NewDirtySet(),
		arrayCtor(t.world),
	}
	worldArray.AddDirtySet(worldArray.dirtySet)
	t.worldArrays[arrayType.String()] = worldArray

	worldCopyArray := entityArray{
		ecs.NewDirtySet(),
		arrayCtor(t.worldCopy),
	}
	t.worldCopyArrays[arrayType.String()] = worldCopyArray

	for _, entity := range worldArray.GetEntities() {
		component, ok := worldArray.GetAny(entity)
		if !ok {
			continue
		}
		worldCopyArray.SetAny(entity, component)
	}
}

func (t tool) SynchronizeState() {
	for arrayType, array := range t.worldArrays {
		t.synchronizeArrayState(
			arrayType,
			array,
			t.worldCopyArrays[arrayType],
		)
	}
}

func (t tool) synchronizeArrayState(
	arrayType string,
	worldArray entityArray,
	worldCopyArray entityArray,
) {
	entities := worldArray.dirtySet.Get()
	if len(entities) == 0 {
		return
	}

	t.applyChangesInEntityRecording(arrayType, worldCopyArray, entities, t.entityKeyedRecorder.recordings, true)
	t.applyChangesInUUIDRecording(arrayType, worldCopyArray, entities, t.uuidKeyedRecorder.uuidRecordings, true)

	for _, entity := range entities {
		if component, ok := worldArray.GetAny(entity); ok {
			worldCopyArray.SetAny(entity, component)
			continue
		}
		if t.world.EntityExists(entity) {
			worldCopyArray.Remove(entity)
			continue
		}
		t.worldCopy.RemoveEntity(entity)
	}

	t.applyChangesInEntityRecording(arrayType, worldCopyArray, entities, t.entityKeyedRecorder.backwardsRecordings, false)
	t.applyChangesInUUIDRecording(arrayType, worldCopyArray, entities, t.uuidKeyedRecorder.backwardsUUIDRecordings, false)
}

func (t tool) applyChangesInEntityRecording(
	arrayType string,
	worldCopyArray ecs.AnyComponentArray,
	entities []ecs.EntityID,
	recordings datastructures.SparseArray[record.RecordingID, Recording],
	seal bool,
) {
	for _, recording := range recordings.GetValues() {
		arrayRecording := recording.Arrays[arrayType]
		for _, entity := range entities {
			if recording.Sealed.Get(entity) {
				continue
			}
			if seal {
				recording.Sealed.Add(entity)
			}

			if component, ok := worldCopyArray.GetAny(entity); ok {
				arrayRecording.Set(entity, component)
				continue
			}
			if t.worldCopy.EntityExists(entity) {
				arrayRecording.Set(entity, nil)
				continue
			}
			recording.RemovedEntities.Add(entity)
		}
	}
}
func (t tool) applyChangesInUUIDRecording(
	arrayType string,
	array ecs.AnyComponentArray,
	entities []ecs.EntityID,
	recordings datastructures.SparseArray[record.UUIDRecordingID, UUIDRecording],
	seal bool,
) {
	for _, recording := range recordings.GetValues() {
		arrayRecording := recording.Arrays[arrayType]
		for _, entity := range entities {
			if recording.Sealed.Get(entity) {
				continue
			}
			if seal {
				recording.Sealed.Add(entity)
			}

			if _, ok := recording.EntitiesUUIDs.Get(entity); !ok {
				uuid, ok := t.world.UUID().Component().Get(entity)
				if !ok {
					uuid.ID = t.world.UUID().NewUUID()
					t.world.UUID().Component().Set(entity, uuid)
				}
				recording.UUIDEntities[uuid.ID] = entity
				recording.EntitiesUUIDs.Set(entity, uuid.ID)
			}

			if component, ok := array.GetAny(entity); ok {
				arrayRecording.Set(entity, component)
				continue
			}
			if t.worldCopy.EntityExists(entity) {
				arrayRecording.Set(entity, nil)
				continue
			}
			recording.RemovedEntities.Add(entity)
		}
	}
}
