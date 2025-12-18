package ecs

import (
	"errors"
	"fmt"
	"reflect"
)

type SystemRegister[World any] interface {
	Register(World) error
}

// impl

type systemRegister[World any] struct{ register func(World) error }

func (s systemRegister[World]) Register(w World) error { return s.register(w) }
func NewSystemRegister[World any](l func(World) error) SystemRegister[World] {
	return systemRegister[World]{l}
}

// helpers

var (
	ErrNotASystem        error = errors.New("system doesn't have proper register method")
	ErrWorldLacksMethods error = errors.New("world doesn't implement system world")
)
var errType = reflect.TypeFor[error]()

func RegisterSystems[World any](w World, systems ...any) []error {
	worldValue := reflect.ValueOf(w)
	worldType := worldValue.Type()
	arguments := []reflect.Value{worldValue}
	errs := []error{}
	for _, system := range systems {
		if system == nil {
			continue
		}
		systemValue := reflect.ValueOf(system)
		register := systemValue.MethodByName("Register")
		if !register.IsValid() {
			errs = append(errs, ErrNotASystem)
			continue
		}
		registerType := register.Type()
		if registerType.NumIn() != 1 ||
			registerType.In(0).Kind() != reflect.Interface ||
			registerType.NumOut() != 1 ||
			!registerType.Out(0).Implements(errType) {
			errs = append(errs, ErrNotASystem)
			continue
		}
		if !worldType.Implements(registerType.In(0)) {
			err := errors.Join(
				ErrWorldLacksMethods,
				fmt.Errorf("not implemented interface type is %v", registerType.In(0).String()),
			)
			errs = append(errs, err)
			continue
		}
		out := register.Call(arguments)
		err, _ := out[0].Interface().(error)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

//

type ToolFactory[World, Tool any] interface {
	Build(World) Tool
}
type toolFactory[World, Tool any] struct{ build func(World) Tool }

func (f toolFactory[World, Tool]) Build(w World) Tool { return f.build(w) }
func NewToolFactory[World, Tool any](l func(World) Tool) ToolFactory[World, Tool] {
	return &toolFactory[World, Tool]{build: l}
}
