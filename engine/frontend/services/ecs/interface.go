package ecs

import (
	"errors"
	"reflect"
)

//
//
// entities

type EntityID struct {
	id string
}

func NewEntityID(id string) EntityID { return EntityID{id} }

func (entityId EntityID) Ok() bool { return entityId.id != "" }

type entitiesInterface interface {
	NewEntity() EntityID
	RemoveEntity(EntityID)

	GetEntities() []EntityID
	EntityExists(EntityID) bool
}

//
//
// components

type ComponentType struct {
	componentType reflect.Type
}

func (t *ComponentType) String() string { return t.componentType.String() }

func newComponentType(componentType reflect.Type) ComponentType {
	return ComponentType{componentType: componentType}
}

type Component interface{}

func GetComponentType(component Component) ComponentType {
	typeOfComponent := reflect.TypeOf(component)
	if typeOfComponent.Kind() != reflect.Struct {
		panic("component has to be a struct (cannot use pointers under the hood)")
	}
	return newComponentType(typeOfComponent)
}

func GetComponentPointerType(componentPointer any) ComponentType {
	return newComponentType(reflect.TypeOf(componentPointer).Elem())
}

var (
	ErrComponentDoNotExists error = errors.New("component do not exists")
	ErrEntityDoNotExists    error = errors.New("entity do not exists")
)

type LiveQuery interface {
	// on add listener add all entities are passed to it
	OnAdd(func([]EntityID))
	OnChange(func([]EntityID))
	OnRemove(func([]EntityID))

	Entities() []EntityID
}

type componentsInterface interface {
	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityID, Component) error // upsert (create or update)
	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetComponent(entityId EntityID, componentType ComponentType) (Component, error)
	RemoveComponent(EntityID, ComponentType)

	// returns for with all listed component types
	// the same live query should be returned for the same input
	QueryEntitiesWithComponents(...ComponentType) LiveQuery
}

func GetComponent[WantedComponent Component](w World, entity EntityID) (WantedComponent, error) {
	var zero WantedComponent
	component, err := w.GetComponent(entity, GetComponentType(zero))
	if err != nil {
		return zero, err
	}
	return component.(WantedComponent), nil
}

//
//
// registry

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

//
//
// world

type World interface {
	entitiesInterface
	componentsInterface
	registryInterface
}
