package ecs

import (
	"errors"
)

type SystemRegister interface {
	Register() error
}

// impl

type systemRegister struct{ register func() error }

func (s systemRegister) Register() error { return s.register() }
func NewSystemRegister(l func() error) SystemRegister {
	return systemRegister{l}
}

// helpers

var (
	ErrNotASystem        error = errors.New("system doesn't have proper register method")
	ErrWorldLacksMethods error = errors.New("world doesn't implement system world")
)

func RegisterSystems(systems ...SystemRegister) []error {
	errs := []error{}
	for _, system := range systems {
		if system == nil {
			continue
		}
		if err := system.Register(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
