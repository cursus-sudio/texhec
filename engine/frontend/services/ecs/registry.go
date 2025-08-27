package ecs

import "sync"

type registryImpl struct {
	registry map[RegisterType]Register
	mutex    *sync.RWMutex
}

func newRegistry(mutex *sync.RWMutex) registryInterface {
	return &registryImpl{
		registry: map[RegisterType]Register{},
		mutex:    mutex,
	}
}

func (r *registryImpl) SaveRegister(register Register) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	registerType := GetRegisterType(register)
	r.registry[registerType] = register
}

func (r *registryImpl) GetRegister(registerType RegisterType) (Register, error) {
	r.mutex.RLocker().Lock()
	defer r.mutex.RLocker().Unlock()
	value, ok := r.registry[registerType]
	if !ok {
		var zero Register
		return zero, ErrRegisterNotFound
	}
	return value, nil
}
