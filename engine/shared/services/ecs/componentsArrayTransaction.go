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
	SaveComponent(EntityID, Component) // upsert
	GetEntityComponent(entity EntityID) EntityComponent[Component]
	AnyComponentsArrayTransaction
}

type AnyComponentsArrayTransaction interface {
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
)

type componentsArrayTransaction[Component any] struct {
	// target of flush
	array *componentsArray[Component]

	// changes
	operations datastructures.SparseArray[EntityID, operation]
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
		saves:      datastructures.NewSparseArray[EntityID, save[Component]](),
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
		case operationRemove:
			t.removes.Remove(entity)
			break
		}
	}
}

func (c *componentsArrayTransaction[Component]) GetEntityComponent(
	entity EntityID,
) EntityComponent[Component] {
	return newEntityComponent(
		entity,
		c.array.GetComponent,
		c.SaveComponent,
		c.RemoveComponent,
	)
}

func (t *componentsArrayTransaction[Component]) SaveComponent(entity EntityID, component Component) {
	t.removeOperation(entity)
	t.saves.Set(entity, save[Component]{entity, component})
	t.operations.Set(entity, operationSave)
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
		value, ok := t.array.components.Get(save.entity)
		if ok && t.array.equal(value, save.component) {
			continue
		}
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

	t.array.applyTransactionMutex.Unlock()

	// notify listeners
	return func() {
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
			for _, listener := range t.array.onRemoveComponents {
				listener(onRemove, onRemoveComponents)
			}
		}
	}, nil
}

func (t *componentsArrayTransaction[Component]) Discard() {
	if t.prepared {
		t.array.applyTransactionMutex.Unlock()
	}
	t.prepared = false
	t.saves = datastructures.NewSparseArray[EntityID, save[Component]]()
	t.removes = datastructures.NewSparseSet[EntityID]()
}

var i int = 0

func FlushMany(transactions ...AnyComponentsArrayTransaction) error {
	i += 1
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
