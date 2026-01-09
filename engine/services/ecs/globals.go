package ecs

import (
	"engine/services/datastructures"
	"reflect"
)

// interface

type GlobalType struct {
	registerType reflect.Type
}

func (t *GlobalType) String() string { return t.registerType.String() }

type Global any

func GetGlobalType(register Global) GlobalType {
	typeOfGlobal := reflect.TypeOf(register)
	if typeOfGlobal.Kind() != reflect.Pointer {
		panic("globals have to be pointers")
	}
	return GlobalType{typeOfGlobal}
}

type globalsInterface interface {
	SaveGlobal(Global) // upsert (create or update)
	GetGlobal(GlobalType) (Global, bool)

	ReleaseGlobals()
}

type Cleanable interface {
	Release()
}

func GetGlobal[GlobalT Global](w World) (*GlobalT, bool) {
	registerType := GetGlobalType((*GlobalT)(nil))
	value, ok := w.GetGlobal(registerType)
	if !ok {
		return nil, ok
	}
	return value.(*GlobalT), true
}

// impl

type globalsImpl struct {
	cleanableTypes datastructures.Set[GlobalType]
	cleanables     datastructures.Array[Cleanable]
	registry       map[GlobalType]Global
}

func newGlobals() *globalsImpl {
	return &globalsImpl{
		cleanableTypes: datastructures.NewSet[GlobalType](),
		cleanables:     datastructures.NewArray[Cleanable](),
		registry:       map[GlobalType]Global{},
	}
}

func (r *globalsImpl) SaveGlobal(register Global) {
	registerType := GetGlobalType(register)
	if cleanable, ok := register.(Cleanable); ok {
		index, ok := r.cleanableTypes.GetIndex(registerType)
		if ok {
			cleanable := r.cleanables.Get()[index]
			cleanable.Release()
			_ = r.cleanables.Remove(index)
			_ = r.cleanableTypes.Remove(index)
		}
		r.cleanableTypes.Add(registerType)
		r.cleanables.Add(cleanable)
	}
	r.registry[registerType] = register
}

func (r *globalsImpl) GetGlobal(registerType GlobalType) (Global, bool) {
	value, ok := r.registry[registerType]
	if !ok {
		return nil, false
	}
	return value, true
}

func (r *globalsImpl) ReleaseGlobals() {
	for _, cleanable := range r.cleanables.Get() {
		cleanable.Release()
	}
	*r = *newGlobals()
}
