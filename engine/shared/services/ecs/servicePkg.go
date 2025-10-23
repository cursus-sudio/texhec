package ecs

import (
	"github.com/ogiusek/events"
)

type SystemRegister interface {
	Register(b events.Builder)
}

func RegisterSystems(b events.Builder, systems ...SystemRegister) {
	for _, system := range systems {
		if system == nil {
			continue
		}
		system.Register(b)
	}
}
