package ecs

import (
	"errors"
	"fmt"
	"frontend/services/datastructures"
	"sync"
)

type registryImpl struct {
	cleanableTypes datastructures.Set[RegisterType]
	cleanables     datastructures.Array[Cleanable]
	registry       map[RegisterType]Register
	mutex          *sync.RWMutex
}

func newRegistry(mutex *sync.RWMutex) *registryImpl {
	return &registryImpl{
		cleanableTypes: datastructures.NewSet[RegisterType](),
		cleanables:     datastructures.NewArray[Cleanable](),
		registry:       map[RegisterType]Register{},
		mutex:          mutex,
	}
}

func (r *registryImpl) SaveRegister(register Register) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	registerType := GetRegisterType(register)
	if cleanable, ok := register.(Cleanable); ok {
		index, ok := r.cleanableTypes.GetIndex(registerType)
		if ok {
			cleanable := r.cleanables.Get()[index]
			cleanable.CleanUp()
			r.cleanables.Remove(index)
			r.cleanableTypes.Remove(index)
		}
		r.cleanableTypes.Add(registerType)
		r.cleanables.Add(cleanable)
	}
	r.registry[registerType] = register
}

func (r *registryImpl) GetRegister(registerType RegisterType) (Register, error) {
	r.mutex.RLocker().Lock()
	defer r.mutex.RLocker().Unlock()
	value, ok := r.registry[registerType]
	if !ok {
		return nil, errors.Join(
			ErrRegisterNotFound,
			fmt.Errorf("haven't found register \"%s\" of type", registerType.String()),
		)
	}
	return value, nil
}

func (r *registryImpl) CleanUp() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, cleanable := range r.cleanables.Get() {
		cleanable.CleanUp()
	}
	*r = *newRegistry(r.mutex)
}
