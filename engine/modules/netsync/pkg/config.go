package netsyncpkg

import (
	"engine/modules/netsync/internal/config"
	"engine/modules/record"
	"reflect"

	"github.com/ogiusek/events"
)

type Config struct {
	config *config.Config
}

func NewConfig(maxPredictions int) Config {
	return Config{
		config: &config.Config{
			RecordConfig:        record.NewConfig(),
			AuthorizeEvent:      make(map[reflect.Type]func(any) error),
			AllowedClientEvents: make(map[reflect.Type]struct{}),
			MaxPredictions:      maxPredictions,
		},
	}
}

func (c Config) RecordConfig() record.Config {
	return c.config.RecordConfig
}

func AddEvent[EventType any](config Config) {
	eventType := reflect.TypeFor[EventType]()
	config.config.EventTypes = append(config.config.EventTypes, eventType)
	config.config.ListenToEvents = append(config.config.ListenToEvents, func(b events.Builder, f func(any)) {
		events.Listen(b, func(e EventType) { f(e) })
	})
	config.config.AllowedClientEvents[eventType] = struct{}{}
}

// these event are sent from server to client regurally but they aren't sent from client to server
func AddSimulatedEvent[EventType any](config Config) {
	config.config.SimulatedEvents = append(config.config.EventTypes, reflect.TypeFor[EventType]())
	config.config.ListenToSimulatedEvents = append(config.config.ListenToEvents, func(b events.Builder, f func(any)) {
		events.Listen(b, func(e EventType) { f(e) })
	})
}

// these are freely exchanged between server and client instead of sending authorized state
func AddTransparentEvent[EventType any](config Config) {
	eventType := reflect.TypeFor[EventType]()
	config.config.TransparentEvents = append(config.config.TransparentEvents, eventType)
	config.config.ListenToTransparentEvents = append(config.config.ListenToTransparentEvents, func(b events.Builder, f func(any)) {
		events.Listen(b, func(e EventType) { f(e) })
	})
	config.config.AllowedClientEvents[eventType] = struct{}{}
}

func AddEventAuthorization[EventType any](config Config, handler func(EventType) error) {
	eventType := reflect.TypeFor[EventType]()
	config.config.AuthorizeEvent[eventType] = func(a any) error {
		return handler(a.(EventType))
	}
}
