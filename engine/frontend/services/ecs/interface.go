package ecs

import (
	"errors"
	"reflect"
)

// entities

type EntityID struct {
	id string
}

func (entityId EntityID) Ok() bool { return entityId.id != "" }

type entitiesInterface interface {
	NewEntity() EntityID
	RemoveEntity(EntityID)

	GetEntities() []EntityID
	EntityExists(EntityID) bool
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
	SaveComponent(EntityID, Component) error // upsert (create or update)
	// can return:
	// - ErrComponentDoNotExists
	// - ErrEntityDoNotExists
	GetComponent(entityId EntityID, componentType ComponentType) (Component, error)
	RemoveComponent(EntityID, ComponentType)

	// returns entities with all listed component types
	GetEntitiesWithComponents(...ComponentType) []EntityID
}

func GetComponent[WantedComponent Component](w World, entity EntityID) (WantedComponent, error) {
	var zero WantedComponent
	component, err := w.GetComponent(entity, GetComponentType(zero))
	if err != nil {
		return zero, err
	}
	return component.(WantedComponent), nil
}

func GetComponents(w World, entity EntityID, wantedComponents ...*Component) error {
	for _, component := range wantedComponents {
		componentType := GetComponentType(*component)
		var err error
		*component, err = w.GetComponent(entity, componentType)
		if err != nil {
			return err
		}
	}
	return nil
}

// world

type World interface {
	entitiesInterface
	componentsInterface
}
