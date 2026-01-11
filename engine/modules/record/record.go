package record

import (
	"engine/modules/uuid"
	"engine/services/datastructures"
	"engine/services/ecs"
	"reflect"
)

type Service interface {
	Entity() EntityKeyedRecorder
	UUID() UUIDKeyedRecorder
}

//

type Config struct {
	ComponentsOrder    *[]reflect.Type
	RecordedComponents map[reflect.Type]func(ecs.World) ecs.AnyComponentArray
}

func NewConfig() Config {
	componentsOrder := make([]reflect.Type, 0)
	return Config{
		ComponentsOrder:    &componentsOrder,
		RecordedComponents: make(map[reflect.Type]func(ecs.World) ecs.AnyComponentArray),
	}
}

func AddToConfig[Component any](config Config) {
	componentType := reflect.TypeFor[Component]()
	if _, ok := config.RecordedComponents[componentType]; ok {
		return
	}
	*config.ComponentsOrder = append(*config.ComponentsOrder, componentType)
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

	Apply(Config, ...Recording)
}

type RecordingID uint16
type Recording struct {
	// [componentArrayLayoutID]component
	// nil for removed entity
	Entities datastructures.SparseArray[ecs.EntityID, []any]
}

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

	Apply(Config, ...UUIDRecording)
}

type UUIDRecordingID uint16
type UUIDRecording struct {
	// map[componentUUID][componentArrayLayoutID]component
	// map[componentUUID]nil is when entity is removed
	Entities map[uuid.UUID][]any
}
