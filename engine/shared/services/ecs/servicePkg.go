package ecs

import (
	"github.com/ogiusek/events"
)

type SystemRegister interface {
	Register(b events.Builder)
}

// impl

type systemRegister struct{ register func(b events.Builder) }

func (s systemRegister) Register(b events.Builder)              { s.register(b) }
func NewSystemRegister(l func(b events.Builder)) SystemRegister { return systemRegister{l} }

// helpers

func RegisterSystems(b events.Builder, systems ...SystemRegister) {
	for _, system := range systems {
		if system == nil {
			continue
		}
		system.Register(b)
	}
}

func ReEmit[EventFrom any](emitter func(e EventFrom)) SystemRegister {
	return NewSystemRegister(func(b events.Builder) {
		events.Listen(b, func(from EventFrom) {
			emitter(from)
		})
	})
}
