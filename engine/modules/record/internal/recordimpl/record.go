package recordimpl

import (
	"engine/modules/record"
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"engine/services/logger"
	"reflect"
	"sync"

	"github.com/ogiusek/ioc/v2"
)

type entityArray struct {
	dirtySet ecs.DirtySet

	// backwards dependencies
	dependencies     datastructures.Set[*BackwardRecording]
	uuidDependencies datastructures.Set[*UUIDBackwardRecording]

	ecs.AnyComponentArray
}

type service struct {
	World       ecs.World    `inject:"1"`
	WorldUUID   uuid.Service `inject:"1"`
	worldArrays map[string]entityArray

	worldCopy       ecs.World
	worldCopyUUID   ecs.ComponentsArray[uuid.Component]
	worldCopyArrays map[string]ecs.AnyComponentArray

	Logger logger.Logger `inject:"1"`
	mutex  *sync.Mutex
	entity *entityKeyedRecorder
	uuid   *uuidKeyedRecorder
}

func NewService(c ioc.Dic) record.Service {
	t := ioc.GetServices[*service](c)
	t.worldArrays = make(map[string]entityArray)

	t.worldCopy = ecs.NewWorld()
	t.worldCopyUUID = ecs.GetComponentsArray[uuid.Component](t.worldCopy)
	t.worldCopyArrays = make(map[string]ecs.AnyComponentArray)

	t.mutex = &sync.Mutex{}

	t.entity = newEntityKeyedRecorder(t)
	t.uuid = newUUIDKeyedRecorder(t)

	return t
}

//

func (t *service) Entity() record.EntityKeyedRecorder {
	return t.entity
}
func (t *service) UUID() record.UUIDKeyedRecorder {
	return t.uuid
}

//

func (t *service) SyncBackwardsRecordingState() {
	for arrayType, array := range t.worldArrays {
		t.synchronizeArrayState(
			array,
			t.worldCopyArrays[arrayType],
		)
	}
}

func (t *service) synchronizeArrayState(
	worldArray entityArray,
	worldCopyArray ecs.AnyComponentArray,
) {
	entities := worldArray.dirtySet.Get()
	if len(entities) == 0 {
		return
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
			uuid, ok := t.worldCopyUUID.Get(entity)
			if !ok {
				uuid, ok = t.WorldUUID.Component().Get(entity)
				if !ok {
					uuid.ID = t.WorldUUID.NewUUID()
					t.WorldUUID.Component().Set(entity, uuid)
				}
				t.worldCopyUUID.Set(entity, uuid)
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

	// apply in world
	for _, entity := range entities {
		if component, ok := worldArray.GetAny(entity); ok {
			err := worldCopyArray.SetAny(entity, component)
			t.Logger.Warn(err)
			continue
		}
		if t.World.EntityExists(entity) {
			worldCopyArray.Remove(entity)
			continue
		}
		t.worldCopy.RemoveEntity(entity)
	}
}

func (t *service) GetWorldArray(arrayType reflect.Type, config record.Config) entityArray {
	arrayKey := arrayType.String()
	if array, ok := t.worldArrays[arrayKey]; ok {
		return array
	}
	arrayCtor := config.RecordedComponents[arrayType]
	entityArray := entityArray{
		dirtySet:          ecs.NewDirtySet(),
		dependencies:      datastructures.NewSet[*BackwardRecording](),
		uuidDependencies:  datastructures.NewSet[*UUIDBackwardRecording](),
		AnyComponentArray: arrayCtor(t.World),
	}
	entityArray.AddDirtySet(entityArray.dirtySet)
	t.worldArrays[arrayKey] = entityArray
	entityArray.dirtySet.Clear()

	array := arrayCtor(t.worldCopy)
	t.worldCopyArrays[arrayKey] = array

	for _, entity := range entityArray.GetEntities() {
		component, _ := entityArray.GetAny(entity)
		_ = array.SetAny(entity, component)
	}

	return entityArray
}

func (t *service) GetWorldCopyArray(arrayType reflect.Type, config record.Config) ecs.AnyComponentArray {
	arrayKey := arrayType.String()
	if array, ok := t.worldCopyArrays[arrayKey]; ok {
		return array
	}
	arrayCtor := config.RecordedComponents[arrayType]
	entityArray := entityArray{
		dirtySet:          ecs.NewDirtySet(),
		dependencies:      datastructures.NewSet[*BackwardRecording](),
		uuidDependencies:  datastructures.NewSet[*UUIDBackwardRecording](),
		AnyComponentArray: arrayCtor(t.World),
	}
	entityArray.AddDirtySet(entityArray.dirtySet)
	t.worldArrays[arrayKey] = entityArray
	entityArray.dirtySet.Clear()

	array := arrayCtor(t.worldCopy)
	t.worldCopyArrays[arrayKey] = array

	for _, entity := range entityArray.GetEntities() {
		component, _ := entityArray.GetAny(entity)
		_ = array.SetAny(entity, component)
	}

	return array
}
