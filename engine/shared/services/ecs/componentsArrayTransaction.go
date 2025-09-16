package ecs

import (
	"errors"
	"fmt"
	"shared/services/datastructures"
)

// interface

type ComponentsArrayTransaction[Component any] interface {
	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) ComponentsArrayTransaction[Component] // upsert
	SaveAnyComponent(EntityID, any) error                                    // upsert
	// differs from save component by not triggering events
	DirtySaveComponent(EntityID, Component) ComponentsArrayTransaction[Component] // upsert
	DirtySaveAnyComponent(EntityID, any) error                                    // upsert

	RemoveComponent(EntityID) ComponentsArrayTransaction[Component]
	RemoveAnyComponent(EntityID)

	// this prematurely locks mutex until we flush (it is optional).
	// it is useful if we want to check error on many arrays.
	PrepareFlush()

	Error() error

	// runs Error() and if it doesn't return an error it applies it.
	// flush also effectively resets transaction to initial state so it can be reused.
	Flush() error
	Discard()
}

type AnyComponentsArrayTransaction interface {
	SaveAnyComponent(EntityID, any) error      // upsert
	DirtySaveAnyComponent(EntityID, any) error // upsert

	RemoveAnyComponent(EntityID)

	// this prematurely locks mutex until we flush (it is optional).
	// it is useful if we want to check error on many arrays.
	PrepareFlush()

	Error() error

	// runs Error() and if it doesn't return an error it applies it.
	// flush also effectively resets transaction to initial state so it can be reused.
	Flush() error
	Discard()
}

// impl

type save[Component any] struct {
	entity    EntityID
	component Component
}

type operation uint8

const (
	operationSave operation = iota
	operationDirtySave
	operationRemove
)

type componentsArrayTransaction[Component any] struct {
	// target of flush
	array *componentsArray[Component]

	// changes
	operations datastructures.SparseArray[EntityID, operation]
	saves      datastructures.SparseArray[EntityID, save[Component]]
	dirtySaves datastructures.SparseArray[EntityID, save[Component]]
	removes    datastructures.SparseSet[EntityID]

	// flush data
	prepared bool
}

func newComponentsArrayTransaction[Component any](
	array *componentsArray[Component],
) ComponentsArrayTransaction[Component] {
	return &componentsArrayTransaction[Component]{
		array: array,

		operations: datastructures.NewSparseArray[EntityID, operation](),
		saves:      datastructures.NewSparseArray[EntityID, save[Component]](),
		dirtySaves: datastructures.NewSparseArray[EntityID, save[Component]](),
		removes:    datastructures.NewSparseSet[EntityID](),

		prepared: false,
	}
}

func (t *componentsArrayTransaction[Component]) removeOperation(entity EntityID) {
	if operation, ok := t.operations.Get(entity); ok {
		t.operations.Remove(entity)
		switch operation {
		case operationSave:
			t.saves.Remove(entity)
			break
		case operationDirtySave:
			t.dirtySaves.Remove(entity)
			break
		case operationRemove:
			t.removes.Remove(entity)
			break
		}
	}
}

func (t *componentsArrayTransaction[Component]) SaveComponent(entity EntityID, component Component) ComponentsArrayTransaction[Component] {
	t.removeOperation(entity)
	t.saves.Set(entity, save[Component]{entity, component})
	t.operations.Set(entity, operationSave)
	return t
}

func (t *componentsArrayTransaction[Component]) SaveAnyComponent(entity EntityID, anyComponent any) error {
	component, ok := anyComponent.(Component)
	if !ok {
		return ErrInvalidType
	}
	t.SaveComponent(entity, component)
	return nil
}

func (t *componentsArrayTransaction[Component]) DirtySaveComponent(entity EntityID, component Component) ComponentsArrayTransaction[Component] {
	t.removeOperation(entity)
	t.dirtySaves.Set(entity, save[Component]{entity, component})
	t.operations.Set(entity, operationDirtySave)
	return t
}

func (t *componentsArrayTransaction[Component]) DirtySaveAnyComponent(entity EntityID, anyComponent any) error {
	component, ok := anyComponent.(Component)
	if !ok {
		return ErrInvalidType
	}
	t.DirtySaveComponent(entity, component)
	return nil
}

func (t *componentsArrayTransaction[Component]) RemoveComponent(entity EntityID) ComponentsArrayTransaction[Component] {
	t.removeOperation(entity)
	t.removes.Add(entity)
	t.operations.Set(entity, operationRemove)
	return t
}

func (t *componentsArrayTransaction[Component]) RemoveAnyComponent(entity EntityID) {
	t.removeOperation(entity)
	t.removes.Add(entity)
	t.operations.Set(entity, operationRemove)
}

func (t *componentsArrayTransaction[Component]) PrepareFlush() {
	t.prepared = true
	t.array.applyTransactionMutex.Lock()
}

func (t *componentsArrayTransaction[Component]) Error() error {
	requiredEntities := make([]EntityID, 0, len(t.saves.GetValues())+len(t.dirtySaves.GetValues()))
	for _, saved := range t.saves.GetValues() {
		requiredEntities = append(requiredEntities, saved.entity)
	}
	for _, saved := range t.dirtySaves.GetValues() {
		requiredEntities = append(requiredEntities, saved.entity)
	}
	for _, entity := range requiredEntities {
		if ok := t.array.entities.Get(entity); !ok {
			return errors.Join(
				ErrEntityDoNotExists,
				fmt.Errorf("missing entity %v", entity),
			)
		}
	}
	return nil
}

// TODO
// currently we don't check for copies due to memory footprint or performance cost.
// best approach is to have one pre-defined array in struct
// (so it doesn't have to be allocated each time).
func (t *componentsArrayTransaction[Component]) Flush() error {
	if !t.prepared {
		t.array.applyTransactionMutex.Lock()
		// unlock happens before listeners
	}
	t.prepared = false

	if err := t.Error(); err != nil {
		t.array.applyTransactionMutex.Unlock()
		return err
	}

	// for listeners
	onAdd := []EntityID{}
	onChange := []EntityID{}
	onRemove := []EntityID{}

	// apply
	t.operations = datastructures.NewSparseArray[EntityID, operation]()
	for _, save := range t.saves.GetValues() {
		added := t.array.components.Set(save.entity, save.component)
		if added {
			onAdd = append(onAdd, save.entity)
		} else {
			onChange = append(onChange, save.entity)
		}
	}
	t.saves = datastructures.NewSparseArray[EntityID, save[Component]]()

	for _, dirtySave := range t.dirtySaves.GetValues() {
		t.array.components.Set(dirtySave.entity, dirtySave.component)
	}
	t.dirtySaves = datastructures.NewSparseArray[EntityID, save[Component]]()

	for _, removedEntity := range t.removes.GetIndices() {
		if removed := t.array.components.Remove(removedEntity); removed {
			onRemove = append(onRemove, removedEntity)
		}
	}
	t.removes = datastructures.NewSparseSet[EntityID]()

	t.array.applyTransactionMutex.Unlock()

	// notify listeners
	if len(onAdd) != 0 {
		for _, listener := range t.array.onAdd {
			listener(onAdd)
		}
	}
	if len(onChange) != 0 {
		for _, listener := range t.array.onChange {
			listener(onChange)
		}
	}
	if len(onRemove) != 0 {
		for _, listener := range t.array.onRemove {
			listener(onRemove)
		}
	}

	return nil
}

func (t *componentsArrayTransaction[Component]) Discard() {
	if !t.prepared {
		t.array.applyTransactionMutex.Lock()
		// unlock happens before listeners
	}
	t.prepared = false
	t.saves = datastructures.NewSparseArray[EntityID, save[Component]]()
	t.dirtySaves = datastructures.NewSparseArray[EntityID, save[Component]]()
	t.removes = datastructures.NewSparseSet[EntityID]()
	t.array.applyTransactionMutex.Unlock()
}

func FlushMany(transactions ...AnyComponentsArrayTransaction) error {
	var err error
	for _, transaction := range transactions {
		transaction.PrepareFlush()
		if err = transaction.Error(); err != nil {
			break
		}
	}
	if err != nil {
		for _, transaction := range transactions {
			transaction.Discard()
		}
		return err
	}

	for _, transaction := range transactions {
		transaction.Flush()
	}
	return nil
}
