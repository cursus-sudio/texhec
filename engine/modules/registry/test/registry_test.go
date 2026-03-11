package test

import (
	"engine/modules/registry"
	"engine/services/ecs"
	"errors"
	"testing"
)

func TestUnusedFieldError(t *testing.T) {
	type TestedStruct struct {
		Field ecs.EntityID `tt:""`
	}
	s := NewSetup()
	instance := TestedStruct{}
	if err := s.Service.Populate(&instance); !errors.Is(err, registry.ErrNotFoundHandlerForAField) {
		t.Errorf("expected err \"%v\" got \"%v\"", registry.ErrNotFoundHandlerForAField, err)
	}
}

func TestUsedField(t *testing.T) {
	type TestedStruct struct {
		Field ecs.EntityID `tag:"value"`
	}
	s := NewSetup()
	instance := TestedStruct{}
	if err := s.Service.Populate(&instance); err != nil {
		t.Errorf("unexpected err \"%v\"", err)
	}
}

func TestWrongInput(t *testing.T) {
	type TestedStruct struct {
		Field ecs.EntityID `tag:"value"`
	}
	s := NewSetup()
	instance := TestedStruct{}
	if err := s.Service.Populate(instance); !errors.Is(err, registry.ErrExpectedPointerToAStruct) {
		t.Errorf("unexpected err \"%v\"", err)
	}
}

func TestDoubleRegister(t *testing.T) {
	s := NewSetup()
	err := s.Service.Register("tag", func(structTagValue string) ecs.EntityID { return 0 })
	if !errors.Is(err, registry.ErrAlreadyRegistered) {
		t.Errorf("expected err \"%v\", got \"%v\"", registry.ErrAlreadyRegistered, err)
	}
}
