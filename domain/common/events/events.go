package events

import "github.com/ogiusek/ioc"

// event errors

type EventErrors interface {
	TopicIsntRegistered(topic string) error
	TopicIsAlreadyRegistered(topic string) error
}

// event

type Event[T any] struct {
	topic   string
	payload T
}

func (event *Event[T]) Topic() string {
	return event.topic
}

func (event *Event[T]) Payload() T {
	return event.payload
}

func NewEvent[T any](topic string, payload T) Event[T] {
	return Event[T]{
		topic:   topic,
		payload: payload,
	}
}

// events

// events:
// have unique handler per topic
// SERVICE
type Events[TInterface any] interface {
	IsTaken(topic string) bool
	RegisterHandler(topic string, handler func(Event[TInterface])) error
	Emit(event Event[TInterface]) error
}

type eventsImpl[TInterface any] struct {
	c        ioc.Dic
	handlers map[string]func(Event[TInterface])
}

func (event *eventsImpl[TInterface]) IsTaken(topic string) bool {
	_, ok := event.handlers[topic]
	return ok
}

func (events *eventsImpl[TInterface]) RegisterHandler(topic string, handler func(Event[TInterface])) error {
	if _, ok := events.handlers[topic]; !ok {
		errors := ioc.Get[EventErrors](events.c)
		return errors.TopicIsAlreadyRegistered(topic)
	}
	events.handlers[topic] = handler
	return nil
}

func (events *eventsImpl[TInterface]) Emit(e Event[TInterface]) error {
	handler, ok := events.handlers[e.topic]
	if !ok {
		return ioc.Get[EventErrors](events.c).TopicIsntRegistered(e.topic)
	}
	handler(e)
	return nil
}

func NewEvents[T any](c ioc.Dic) Events[T] {
	return &eventsImpl[T]{
		c:        c,
		handlers: make(map[string]func(Event[T])),
	}
}
