package config

import (
	"engine/services/ecs"
	"reflect"

	"github.com/ogiusek/events"
)

type Config struct {
	Events         []reflect.Type
	ListenToEvents []func(events.Builder, func(any))

	TransparentEvents         []reflect.Type
	ListenToTransparentEvents []func(events.Builder, func(any))

	Components         []reflect.Type
	ArraysOfComponents []func(ecs.World) ecs.AnyComponentArray

	// client
	MaxPredictions int

	// server
}
