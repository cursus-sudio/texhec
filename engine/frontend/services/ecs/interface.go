package ecs

import (
	"errors"
	"reflect"
)

// entities

type EntityId struct {
	id string
}

func (entityId EntityId) Ok() bool { return entityId.id != "" }

type entitiesInterface interface {
	NewEntity() EntityId
	RemoveEntity(EntityId)

	GetEntities() []EntityId
	EntityExists(EntityId) bool
}

// components

type ComponentType struct {
	componentType reflect.Type
}

func NewComponentType(componentType reflect.Type) ComponentType {
	return ComponentType{componentType: componentType}
}

type Component interface{}

func GetComponentType(component Component) ComponentType {
	typeOfComponent := reflect.TypeOf(component)
	if typeOfComponent.Kind() != reflect.Struct {
		panic("component has to be a struct")
	}
	return NewComponentType(reflect.TypeOf(component))
}

func GetComponentPointerType(componentPointer any) ComponentType {
	return NewComponentType(reflect.TypeOf(componentPointer).Elem())
}

var (
	ErrComponentDoNotExists error = errors.New("component do not exists")
	ErrEntityDoNotExists    error = errors.New("entity do not exists")
)

type componentsInterface interface {
	// can return:
	// - ErrEntityDoNotExists
	SaveComponent(EntityId, Component) error // upsert (create or update)
	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetComponents(entityId EntityId, componentsPointers ...any) error
	RemoveComponent(EntityId, ComponentType)

	// returns entities with all listed component types
	GetEntitiesWithComponents(...ComponentType) []EntityId
}

// world

type World interface {
	entitiesInterface
	componentsInterface
}
