package ecs

import (
	"errors"
	"fmt"
	"frontend/services/datastructures"
	"reflect"
	"sync"
)

// interface

var (
	ErrRegisterNotFound error = errors.New("register not found")
)

type RegisterType struct {
	registerType reflect.Type
}

func (t *RegisterType) String() string { return t.registerType.String() }

type Register any

func GetRegisterType(register Register) RegisterType {
	typeOfRegister := reflect.TypeOf(register)
	if typeOfRegister.Kind() != reflect.Struct {
		panic("register has to be a struct (can use pointers under the hood)")
	}
	return RegisterType{typeOfRegister}
}

type registryInterface interface {
	SaveRegister(Register) // upsert (create or update)
	GetRegister(RegisterType) (Register, error)

	Release()

	LockRegistry()
	UnlockRegistry()
}

type Cleanable interface {
	Release()
}

func GetRegister[RegisterT Register](w World) (RegisterT, error) {
	var zero RegisterT
	registerType := GetRegisterType(zero)
	value, err := w.GetRegister(registerType)
	if err != nil {
		return zero, err
	}
	return value.(RegisterT), nil
}

// impl

type registryImpl struct {
	cleanableTypes datastructures.Set[RegisterType]
	cleanables     datastructures.Array[Cleanable]
	registry       map[RegisterType]Register
	mutex          sync.Locker
}

func newRegistry() *registryImpl {
	return &registryImpl{
		cleanableTypes: datastructures.NewSet[RegisterType](),
		cleanables:     datastructures.NewArray[Cleanable](),
		registry:       map[RegisterType]Register{},
		mutex:          &sync.Mutex{},
	}
}

func (r *registryImpl) SaveRegister(register Register) {
	registerType := GetRegisterType(register)
	if cleanable, ok := register.(Cleanable); ok {
		index, ok := r.cleanableTypes.GetIndex(registerType)
		if ok {
			cleanable := r.cleanables.Get()[index]
			cleanable.Release()
			r.cleanables.Remove(index)
			r.cleanableTypes.Remove(index)
		}
		r.cleanableTypes.Add(registerType)
		r.cleanables.Add(cleanable)
	}
	r.registry[registerType] = register
}

func (r *registryImpl) GetRegister(registerType RegisterType) (Register, error) {
	value, ok := r.registry[registerType]
	if !ok {
		return nil, errors.Join(
			ErrRegisterNotFound,
			fmt.Errorf("haven't found register \"%s\" of type", registerType.String()),
		)
	}
	return value, nil
}

func (r *registryImpl) Release() {
	for _, cleanable := range r.cleanables.Get() {
		cleanable.Release()
	}
	*r = *newRegistry()
}

func (r *registryImpl) LockRegistry()   { r.mutex.Lock() }
func (r *registryImpl) UnlockRegistry() { r.mutex.Unlock() }
