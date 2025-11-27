package ecs

import (
	"engine/services/datastructures"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// interface

var (
	ErrGlobalNotFound error = errors.New("globl not found")
)

type GlobalType struct {
	registerType reflect.Type
}

func (t *GlobalType) String() string { return t.registerType.String() }

type Global any

func GetGlobalType(register Global) GlobalType {
	typeOfGlobal := reflect.TypeOf(register)
	if typeOfGlobal.Kind() != reflect.Struct {
		panic("register has to be a struct (can use pointers under the hood)")
	}
	return GlobalType{typeOfGlobal}
}

type globalsInterface interface {
	SaveGlobal(Global) // upsert (create or update)
	GetGlobal(GlobalType) (Global, error)

	Release()

	LockGlobals()
	UnlockGlobals()
}

type Cleanable interface {
	Release()
}

func GetGlobal[GlobalT Global](w World) (GlobalT, error) {
	var zero GlobalT
	registerType := GetGlobalType(zero)
	value, err := w.GetGlobal(registerType)
	if err != nil {
		return zero, err
	}
	return value.(GlobalT), nil
}

// impl

type globalsImpl struct {
	cleanableTypes datastructures.Set[GlobalType]
	cleanables     datastructures.Array[Cleanable]
	registry       map[GlobalType]Global
	mutex          sync.Locker
}

func newGlobals() *globalsImpl {
	return &globalsImpl{
		cleanableTypes: datastructures.NewSet[GlobalType](),
		cleanables:     datastructures.NewArray[Cleanable](),
		registry:       map[GlobalType]Global{},
		mutex:          &sync.Mutex{},
	}
}

func (r *globalsImpl) SaveGlobal(register Global) {
	registerType := GetGlobalType(register)
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

func (r *globalsImpl) GetGlobal(registerType GlobalType) (Global, error) {
	value, ok := r.registry[registerType]
	if !ok {
		return nil, errors.Join(
			ErrGlobalNotFound,
			fmt.Errorf("haven't found global \"%s\" of type", registerType.String()),
		)
	}
	return value, nil
}

func (r *globalsImpl) Release() {
	for _, cleanable := range r.cleanables.Get() {
		cleanable.Release()
	}
	*r = *newGlobals()
}

func (r *globalsImpl) LockGlobals()   { r.mutex.Lock() }
func (r *globalsImpl) UnlockGlobals() { r.mutex.Unlock() }
