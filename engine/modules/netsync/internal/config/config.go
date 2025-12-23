package config

import (
	"engine/modules/netsync"
	"engine/modules/record"
	"engine/services/ecs"
	"reflect"

	"github.com/ogiusek/events"
)

type Config struct {
	EventTypes     []reflect.Type
	ListenToEvents []func(events.Builder, func(any))

	SimulatedEvents         []reflect.Type
	ListenToSimulatedEvents []func(events.Builder, func(any))

	TransparentEvents         []reflect.Type
	ListenToTransparentEvents []func(events.Builder, func(any))

	RecordConfig record.Config

	// client
	MaxPredictions int

	// auth
	AllowedClientEvents map[reflect.Type]struct{}
	AuthorizeEvent      map[reflect.Type]func(any) error
}

func (config Config) Auth(client ecs.EntityID, event any) (any, error) {
	eventValue := reflect.ValueOf(event)
	eventType := eventValue.Type()
	if _, ok := config.AllowedClientEvents[eventType]; !ok {
		return event, nil
	}
	eventPointerValue := reflect.New(eventType)
	eventPointerValue.Elem().Set(eventValue)

	eventPointer := eventPointerValue.Interface()
	if authorizedEvent, ok := eventPointer.(netsync.AuthorizedEvent); ok {
		authorizedEvent.SetConnection(client)
	}

	event = eventPointerValue.Elem().Interface()
	if handler, ok := config.AuthorizeEvent[eventType]; ok {
		if err := handler(event); err != nil {
			return nil, err
		}
	}
	return event, nil
}
