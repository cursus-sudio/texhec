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

	// backwards dependencies
	dependencies     datastructures.Set[*BackwardRecording]
	uuidDependencies datastructures.Set[*UUIDBackwardRecording]

	ecs.AnyComponentArray
}

type tool struct {
	world       record.World
	worldArrays map[string]entityArray

	worldCopy       record.World
	worldCopyArrays map[string]ecs.AnyComponentArray

	logger logger.Logger
	mutex  *sync.Mutex
	entity *entityKeyedRecorder
	uuid   *uuidKeyedRecorder
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
			worldCopyArrays: make(map[string]ecs.AnyComponentArray),

			logger: logger,
			mutex:  &sync.Mutex{},
		}
		t.entity = newEntityKeyedRecorder(&t)
		t.uuid = newUUIDKeyedRecorder(&t)
		w.SaveGlobal(t)

		return t
	})
}

func (t tool) Record() record.Interface {
	return t
}

//

func (t tool) Entity() record.EntityKeyedRecorder {
	return t.entity
}
func (t tool) UUID() record.UUIDKeyedRecorder {
	return t.uuid
}

//

func (t tool) SyncBackwardsRecordingState() {
	for arrayType, array := range t.worldArrays {
		t.synchronizeArrayState(
			array,
			t.worldCopyArrays[arrayType],
		)
	}
}

func (t tool) synchronizeArrayState(
	worldArray entityArray,
	worldCopyArray ecs.AnyComponentArray,
) {
	entities := worldArray.dirtySet.Get()
	if len(entities) == 0 {
		return
	}

	// apply in world
	for _, entity := range entities {
		if component, ok := worldArray.GetAny(entity); ok {
			err := worldCopyArray.SetAny(entity, component)
			t.logger.Warn(err)
			continue
		}
		if t.world.EntityExists(entity) {
			worldCopyArray.Remove(entity)
			continue
		}
		t.worldCopy.RemoveEntity(entity)
	}

	// apply in Entity arrays
	for _, recording := range t.entity.backwardsRecordings.GetValues() {
		for _, entity := range entities {
			if _, ok := recording.Entities.Get(entity); ok {
				continue
			}
			var components []any
			if !t.worldCopy.EntityExists(entity) {
				goto saveEntity
			}
			components = make([]any, 0, len(recording.WorldCopyArrays))
			for _, array := range recording.WorldCopyArrays {
				component, ok := array.GetAny(entity)
				if !ok {
					component = nil
				}
				components = append(components, component)
			}
		saveEntity:
			recording.Entities.Set(entity, components)
		}
	}

	// apply in UUID arrays
	for _, recording := range t.uuid.backwardsRecordings.GetValues() {
		for _, entity := range entities {
			uuid, ok := t.worldCopy.UUID().Component().Get(entity)
			if !ok {
				uuid, ok = t.world.UUID().Component().Get(entity)
				if !ok {
					uuid.ID = t.world.UUID().NewUUID()
					t.world.UUID().Component().Set(entity, uuid)
				}
				t.worldCopy.UUID().Component().Set(entity, uuid)
			}
			if _, ok := recording.Entities[uuid.ID]; ok {
				continue
			}
			var components []any
			if !t.worldCopy.EntityExists(entity) {
				goto saveUUID
			}
			components = make([]any, 0, len(recording.WorldCopyArrays))
			for _, array := range recording.WorldCopyArrays {
				component, ok := array.GetAny(entity)
				if !ok {
					component = nil
				}
				components = append(components, component)
			}
		saveUUID:
			recording.Entities[uuid.ID] = components
		}
	}
}

func (t tool) GetWorldArray(arrayType reflect.Type, config record.Config) entityArray {
	arrayKey := arrayType.String()
	if array, ok := t.worldArrays[arrayKey]; ok {
		return array
	}
	arrayCtor := config.RecordedComponents[arrayType]
	entityArray := entityArray{
		dirtySet:          ecs.NewDirtySet(),
		dependencies:      datastructures.NewSet[*BackwardRecording](),
		uuidDependencies:  datastructures.NewSet[*UUIDBackwardRecording](),
		AnyComponentArray: arrayCtor(t.world),
	}
	entityArray.AddDirtySet(entityArray.dirtySet)
	t.worldArrays[arrayKey] = entityArray
	array := arrayCtor(t.worldCopy)
	t.worldCopyArrays[arrayKey] = array
	return entityArray
}

func (t tool) GetWorldCopyArray(arrayType reflect.Type, config record.Config) ecs.AnyComponentArray {
	arrayKey := arrayType.String()
	if array, ok := t.worldCopyArrays[arrayKey]; ok {
		return array
	}
	arrayCtor := config.RecordedComponents[arrayType]
	entityArray := entityArray{
		dirtySet:          ecs.NewDirtySet(),
		dependencies:      datastructures.NewSet[*BackwardRecording](),
		uuidDependencies:  datastructures.NewSet[*UUIDBackwardRecording](),
		AnyComponentArray: arrayCtor(t.world),
	}
	entityArray.AddDirtySet(entityArray.dirtySet)
	t.worldArrays[arrayKey] = entityArray
	array := arrayCtor(t.worldCopy)
	t.worldCopyArrays[arrayKey] = array
	return array
}
