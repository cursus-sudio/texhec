package record

import (
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"reflect"
)

type ToolFactory ecs.ToolFactory[World, RecordTool]
type RecordTool interface {
	Record() Interface
}
type World interface {
	ecs.World
	uuid.UUIDTool
}
type Interface interface {
	Entity() EntityKeyedRecorder
	UUID() UUIDKeyedRecorder
}

//

type Config struct {
	RecordedComponents map[reflect.Type]func(ecs.World) ecs.AnyComponentArray
}

func NewConfig() Config {
	return Config{
		RecordedComponents: make(map[reflect.Type]func(ecs.World) ecs.AnyComponentArray),
	}
}

func AddToConfig[Component any](config Config) {
	componentType := reflect.TypeFor[Component]()
	config.RecordedComponents[componentType] = func(w ecs.World) ecs.AnyComponentArray {
		return ecs.GetComponentsArray[Component](w)
	}
}

//

type EntityKeyedRecorder interface {
	// gets state as finished recording
	GetState(Config) Recording

	// starts opened recording (opened recording is recorded until stopped)
	// applying it on previous state will create current state
	StartRecording(Config) RecordingID
	// starts opened recording (opened recording is recorded until stopped)
	// applying it rewinds state.
	StartBackwardsRecording(Config) RecordingID
	// finishes recording if open
	Stop(RecordingID) (r Recording, ok bool)

	Apply(...Recording)
}

type RecordingID uint16
type Recording struct {
	RemovedEntities datastructures.SparseSet[ecs.EntityID]
	Sealed          datastructures.SparseSet[ecs.EntityID]
	Arrays          map[reflect.Type]ArrayRecording
}
type ArrayRecording datastructures.SparseArray[ecs.EntityID, any] // nil for component means that component is removed

//

type UUIDKeyedRecorder interface {
	// gets state as finished recording
	GetState(Config) UUIDRecording

	// starts opened recording (opened recording is recorded until stopped)
	// applying it on previous state will create current state
	StartRecording(Config) UUIDRecordingID
	// starts opened recording (opened recording is recorded until stopped)
	// applying it rewinds state.
	StartBackwardsRecording(Config) UUIDRecordingID
	// finishes recording if open
	Stop(UUIDRecordingID) (r UUIDRecording, ok bool)

	Apply(...UUIDRecording)
}

type UUIDRecordingID uint16
type UUIDRecording struct {
	UUIDEntities  map[uuid.UUID]ecs.EntityID
	EntitiesUUIDs datastructures.SparseArray[ecs.EntityID, uuid.UUID]
	Recording
}
