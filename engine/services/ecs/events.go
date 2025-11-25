package ecs

import "github.com/ogiusek/events"

type eventsInterface interface {
	Events() events.Events
	EventsBuilder() events.Builder
}

type eventsImpl struct {
	events        events.Events
	eventsBuilder events.Builder
}

func (i *eventsImpl) Events() events.Events         { return i.events }
func (i *eventsImpl) EventsBuilder() events.Builder { return i.eventsBuilder }

func newEvents(b events.Builder) eventsInterface {
	return &eventsImpl{
		events:        b.Events(),
		eventsBuilder: b,
	}
}
