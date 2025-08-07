package colliders

import (
	"errors"
	"reflect"
)

var (
	ErrCollidersHaveDifferentLayersCount error = errors.New("colliders layers count doesn't match")
	ErrMissingShapeHandler               error = errors.New("missing shape handler")
)

type shapeHandlerID struct{ s1, s2 reflect.Type }
type shapeHandler func(Shape, Shape) Collision

type ColliderService interface {
	// collision can be nil if its nil then there is no collision
	Collides(c1, c2 Collider) (Collision, error)
}

type colliderService struct {
	shapeHandlers map[shapeHandlerID]shapeHandler
}

func (colliderService *colliderService) Collides(c1, c2 Collider) (Collision, error) {
	c1Layers := c1.Layers
	c2Layers := c2.Layers
	layersCount := len(c1Layers)
	if len(c2Layers) != layersCount {
		return nil, ErrCollidersHaveDifferentLayersCount
	}
	for i := 0; i < layersCount-1; i++ {
		collides := false
		l1, l2 := c1Layers[i], c2Layers[i]
		for _, s1 := range l1 {
			for _, s2 := range l2 {
				handlerID := shapeHandlerID{s1: reflect.TypeOf(s1), s2: reflect.TypeOf(s2)}
				handler, ok := colliderService.shapeHandlers[handlerID]
				if !ok {
					return nil, ErrMissingShapeHandler
				}
				intersection := handler(s1, s2)
				collides = intersection != nil
				if collides {
					break
				}
			}
			if collides {
				break
			}
		}
		if collides {
			continue
		}
		return nil, nil
	}
	{
		i := layersCount - 1
		l1, l2 := c1Layers[i], c2Layers[i]
		collides := false
		intersections := []Intersection{}
		for _, s1 := range l1 {
			for _, s2 := range l2 {
				handlerID := shapeHandlerID{s1: reflect.TypeOf(s1), s2: reflect.TypeOf(s2)}
				handler, ok := colliderService.shapeHandlers[handlerID]
				if !ok {
					return nil, ErrMissingShapeHandler
				}
				collision := handler(s1, s2)
				if collision == nil {
					continue
				}
				collides = true
				intersections = append(intersections, collision.Intersections()...)
			}
		}
		if !collides {
			return nil, nil
		}
		return NewCollision(intersections...), nil
	}
}
