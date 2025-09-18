package saves

import (
	"encoding/json"
	"errors"
	"reflect"
	"shared/services/codec"
	"shared/services/ecs"
	"sync"
)

type WorldStateCodecBuilder struct {
	arrays map[RepoId]func(ecs.World) ecs.AnyComponentArray
}

func NewWorldStateCodecBuilder() WorldStateCodecBuilder {
	return WorldStateCodecBuilder{
		arrays: make(map[RepoId]func(ecs.World) ecs.AnyComponentArray),
	}
}

func AddPersistedArray[ComponentType any](b WorldStateCodecBuilder) {
	repoName := reflect.TypeFor[ComponentType]().String()
	repoID := RepoId(repoName)
	b.arrays[repoID] = func(w ecs.World) ecs.AnyComponentArray {
		return ecs.GetComponentsArray[ComponentType](w.Components())
	}
}

type worldStateCodec struct {
	repoWMutex sync.Locker
	world      ecs.World
	codec      codec.Codec
	arrays     map[RepoId]func(ecs.World) ecs.AnyComponentArray
}

func newStateCodec(
	repoWMutex sync.Locker,
	world ecs.World,
	codec codec.Codec,
	arrays map[RepoId]func(ecs.World) ecs.AnyComponentArray,
) StateCodec {
	return &worldStateCodec{
		repoWMutex: repoWMutex,
		world:      world,
		codec:      codec,
		arrays:     arrays,
	}
}

type serializableElement struct {
	Entity    ecs.EntityID
	Component []byte
}

func (repoStateCodec *worldStateCodec) Serialize() SaveData {
	repoStateCodec.repoWMutex.Lock()
	defer repoStateCodec.repoWMutex.Unlock()
	serializable := make(map[RepoId][]serializableElement, len(repoStateCodec.arrays))
	for repoId, getter := range repoStateCodec.arrays {
		array := getter(repoStateCodec.world)
		entities := array.GetEntities()
		data := make([]serializableElement, 0, len(entities))
		for _, entity := range entities {
			component, _ := array.GetAnyComponent(entity)
			componentSerialized := repoStateCodec.codec.Encode(component)
			serialzied := serializableElement{entity, componentSerialized}
			data = append(data, serialzied)
		}
		serializable[repoId] = data
	}
	bytes, _ := json.Marshal(serializable)
	return NewSaveData(bytes)
}

func (repoStateCodec *worldStateCodec) Load(data SaveData) error {
	repoStateCodec.repoWMutex.Lock()
	defer repoStateCodec.repoWMutex.Unlock()

	// get repositories snapshots
	snapshots := make(map[RepoId][]serializableElement, len(repoStateCodec.arrays))
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return ErrInvalidSaveFormat
	}

	transactions := make(map[RepoId]ecs.AnyComponentsArrayTransaction, len(repoStateCodec.arrays))

	for key := range snapshots {
		getter, ok := repoStateCodec.arrays[key]
		if !ok {
			return errors.Join(
				ErrInvalidRepoSnapshot,
				errors.New("repo isn't version compatible"),
			)
		}
		array := getter(repoStateCodec.world)
		transaction := array.AnyTransaction()
		transactions[key] = transaction
	}

	// apply snapshots
	for key, snapshot := range snapshots {
		getter := repoStateCodec.arrays[key]
		array := getter(repoStateCodec.world)
		transaction := transactions[key]
		for _, entity := range array.GetEntities() {
			transaction.RemoveAnyComponent(entity)
		}

		for _, element := range snapshot {
			encodedComponent := element.Component
			decodedComponent, err := repoStateCodec.codec.Decode(encodedComponent)
			if err != nil {
				return ErrInvalidRepoSnapshot
			}
			if err := transaction.SaveAnyComponent(element.Entity, decodedComponent); err != nil {
				return errors.Join(ErrInvalidRepoSnapshot, err)
			}
		}
	}
	transactionsArray := make([]ecs.AnyComponentsArrayTransaction, 0, len(transactions))
	for _, transaction := range transactions {
		transactionsArray = append(transactionsArray, transaction)
	}
	return ecs.FlushMany(transactionsArray...)
}
