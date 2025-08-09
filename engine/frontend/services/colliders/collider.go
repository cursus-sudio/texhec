package colliders

import "frontend/engine/components/transform"

type ColliderLayer []Shape

func (originalLayer ColliderLayer) Apply(transform transform.Transform) ColliderLayer {
	layer := make([]Shape, 0, len(originalLayer))
	for _, shape := range originalLayer {
		layer = append(layer, shape.Apply(transform))
	}
	return layer
}

type Collider struct{ collider }

type collider struct {
	Layers []ColliderLayer
}

func NewCollider(layers ...ColliderLayer) Collider {
	return Collider{collider{Layers: layers}}
}

func (collider collider) Apply(transform transform.Transform) Collider {
	layers := make([]ColliderLayer, 0, len(collider.Layers))
	for _, layer := range collider.Layers {
		layers = append(layers, layer.Apply(transform))
	}
	return NewCollider(layers...)
}
