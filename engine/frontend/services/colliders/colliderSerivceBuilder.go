package colliders

import (
	"errors"
	"reflect"
)

var ErrAlreadyRegisteredHandler error = errors.New("handler is already registered")

type ColliderServiceBuilder interface {
	AddHandler(s1, s2 reflect.Type, handler shapeHandler)
	Build() (ColliderService, []error)
}

type colliderServiceBuilder struct {
	handlers map[shapeHandlerID]shapeHandler
	errors   []error
}

func NewBuilder() ColliderServiceBuilder {
	return &colliderServiceBuilder{
		handlers: make(map[shapeHandlerID]shapeHandler),
	}
}

func (c *colliderServiceBuilder) AddHandler(s1, s2 reflect.Type, handler shapeHandler) {
	id := shapeHandlerID{s1: s1, s2: s2}
	if _, ok := c.handlers[id]; ok {
		c.errors = append(c.errors, ErrAlreadyRegisteredHandler)
		return
	}
	c.handlers[id] = handler
	if s1 == s2 {
		return
	}
	c.handlers[shapeHandlerID{s1: s2, s2: s1}] = func(s1, s2 Shape) Collision {
		collision := handler(s2, s1)
		if collision == nil {
			return nil
		}
		return collision.Reverse()
	}
}

func AddHandler[S1, S2 Shape](b ColliderServiceBuilder, handler func(s1 S1, s2 S2) Collision) {
	t1, t2 := reflect.TypeFor[S1](), reflect.TypeFor[S2]()
	b.AddHandler(t1, t2, func(s1, s2 Shape) Collision {
		return handler(s1.(S1), s2.(S2))
	})
}

func (c *colliderServiceBuilder) Build() (ColliderService, []error) {
	if len(c.errors) != 0 {
		return nil, c.errors
	}
	return &colliderService{shapeHandlers: c.handlers}, nil
}
