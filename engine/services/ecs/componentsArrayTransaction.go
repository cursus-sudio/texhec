package ecs

import (
	"engine/services/datastructures"
	"errors"
	"fmt"
)

// interface

type ComponentsArrayTransaction[Component any] interface {
	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) // upsert
	GetEntityComponent(entity EntityID) EntityComponent[Component]
	AnyComponentsArrayTransaction
}

type AnyComponentsArrayTransaction interface {
	TriggerChangeListener(EntityID)
	SaveAnyComponent(EntityID, any) error // upsert

	RemoveComponent(EntityID)

	// this prematurely locks mutex until we flush (it is optional).
	// it is useful if we want to check error on many arrays.
	PrepareFlush()

	Error() error

	// runs Error() and if it doesn't return an error it applies it.
	// flush also effectively resets transaction to initial state so it can be reused.
	Flush() (listeners func(), err error)
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
	operationRemove
	operationChanged
)

type componentsArrayTransaction[Component any] struct {
	// target of flush
	array *componentsArray[Component]

	// changes
	operations datastructures.SparseArray[EntityID, operation]
	changes    datastructures.SparseSet[EntityID]
	saves      datastructures.SparseArray[EntityID, save[Component]]
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
		changes:    datastructures.NewSparseSet[EntityID](),
		saves:      datastructures.NewSparseArray[EntityID, save[Component]](),
		removes:    datastructures.NewSparseSet[EntityID](),
		prepared:   false,
	}
}

func (t *componentsArrayTransaction[Component]) removeOperation(entity EntityID) {
	if operation, ok := t.operations.Get(entity); ok {
		t.operations.Remove(entity)
		switch operation {
		case operationSave:
			t.saves.Remove(entity)
		case operationRemove:
			t.removes.Remove(entity)
		case operationChanged:
			t.changes.Remove(entity)
		}
	}
}

func (t *componentsArrayTransaction[Component]) GetEntityComponent(
	entity EntityID,
) EntityComponent[Component] {
	return newEntityComponent(
		entity,
		t.array.GetComponent,
		t.SaveComponent,
		t.RemoveComponent,
	)
}

func (t *componentsArrayTransaction[Component]) SaveComponent(entity EntityID, component Component) {
	comp, err := t.array.GetComponent(entity)
	if err == nil && t.array.equal(comp, component) {
		return
	}
	t.removeOperation(entity)
	t.saves.Set(entity, save[Component]{entity, component})
	t.operations.Set(entity, operationSave)
}

func (t *componentsArrayTransaction[Component]) TriggerChangeListener(entity EntityID) {
	if _, ok := t.operations.Get(entity); !ok {
		t.operations.Set(entity, operationChanged)
		t.changes.Add(entity)
	}
}

func (t *componentsArrayTransaction[Component]) SaveAnyComponent(entity EntityID, anyComponent any) error {
	component, ok := anyComponent.(Component)
	if !ok {
		return ErrInvalidType
	}
	t.SaveComponent(entity, component)
	return nil
}

func (t *componentsArrayTransaction[Component]) RemoveComponent(entity EntityID) {
	t.removeOperation(entity)
	t.removes.Add(entity)
	t.operations.Set(entity, operationRemove)
}

func (t *componentsArrayTransaction[Component]) PrepareFlush() {
	if t.prepared {
		return
	}
	t.prepared = true
	t.array.applyTransactionMutex.Lock()
}

func (t *componentsArrayTransaction[Component]) Error() error {
	requiredEntities := make([]EntityID, 0, len(t.saves.GetValues()))
	for _, saved := range t.saves.GetValues() {
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

func (t *componentsArrayTransaction[Component]) Flush() (func(), error) {
	if !t.prepared {
		t.array.applyTransactionMutex.Lock()
		// unlock happens before listeners
	}
	t.prepared = false

	if len(t.operations.GetIndices()) == 0 {
		t.array.applyTransactionMutex.Unlock()
		return func() {}, nil
	}

	if err := t.Error(); err != nil {
		t.array.applyTransactionMutex.Unlock()
		return nil, err
	}

	// for listeners
	onAdd := []EntityID{}
	onChange := []EntityID{}
	onRemove := []EntityID{}
	onRemoveComponents := []Component{}

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

	for _, removedEntity := range t.removes.GetIndices() {
		component, _ := t.array.components.Get(removedEntity)
		if removed := t.array.components.Remove(removedEntity); removed {
			onRemove = append(onRemove, removedEntity)
			onRemoveComponents = append(onRemoveComponents, component)
		}
	}
	t.removes = datastructures.NewSparseSet[EntityID]()

	for _, entity := range t.changes.GetIndices() {
		onChange = append(onChange, entity)
	}
	t.changes = datastructures.NewSparseSet[EntityID]()

	t.array.applyTransactionMutex.Unlock()

	// notify listeners
	return func() {
		addI := 0
		changeI := 0
		removeI := 0
		removeComponentsI := 0
		for _, listener := range t.array.listenersOrder {
			switch listener {
			case addListener:
				if len(onAdd) != 0 {
					t.array.onAdd[addI](onAdd)
				}
				addI++
			case changeListener:
				if len(onChange) != 0 {
					t.array.onChange[changeI](onChange)
				}
				changeI++
			case removeListener:
				if len(onRemove) != 0 {
					t.array.onRemove[removeI](onRemove)
				}
				removeI++
			case removeComponentsListener:
				if len(onRemoveComponents) != 0 {
					t.array.onRemoveComponents[removeComponentsI](onRemove, onRemoveComponents)
				}
				removeComponentsI++
			}
		}
		// if len(onAdd) != 0 {
		// 	for _, listener := range t.array.onAdd {
		// 		listener(onAdd)
		// 	}
		// }
		// if len(onChange) != 0 {
		// 	for _, listener := range t.array.onChange {
		// 		listener(onChange)
		// 	}
		// }
		// if len(onRemove) != 0 {
		// 	for _, listener := range t.array.onRemove {
		// 		listener(onRemove)
		// 	}
		// 	for _, listener := range t.array.onRemoveComponents {
		// 		listener(onRemove, onRemoveComponents)
		// 	}
		// }
	}, nil
}

func (t *componentsArrayTransaction[Component]) Discard() {
	if t.prepared {
		t.array.applyTransactionMutex.Unlock()
	}
	t.prepared = false
	t.operations = datastructures.NewSparseArray[EntityID, operation]()
	t.saves = datastructures.NewSparseArray[EntityID, save[Component]]()
	t.changes = datastructures.NewSparseSet[EntityID]()
	t.removes = datastructures.NewSparseSet[EntityID]()
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

	listeners := make([]func(), 0, len(transactions))
	for _, transaction := range transactions {
		listener, _ := transaction.Flush()
		listeners = append(listeners, listener)
	}

	for _, listener := range listeners {
		listener()
	}
	return nil
}
