package colliders

type ColliderLayer []Shape

type Collider interface {
	Layers() []ColliderLayer
}

type collider struct {
	layers []ColliderLayer
}

func (collider *collider) Layers() []ColliderLayer { return collider.layers }

func NewCollider(layers ...ColliderLayer) Collider {
	return &collider{layers: layers}
}
