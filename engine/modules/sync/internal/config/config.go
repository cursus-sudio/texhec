package config

import (
	"engine/services/ecs"
	"reflect"

	"github.com/ogiusek/events"
)

type Config struct {
	// shared
	// events and their usages
	Events         []reflect.Type
	ListenToEvents []func(events.Builder, func(any))

	// components and their usages
	Components         []reflect.Type
	ArraysOfComponents []func(ecs.World) ecs.AnyComponentArray

	IsClient bool

	// client
	MaxPredictions int

	// server
}
