package colliders

type ColliderLayer []Shape

type Collider struct{ collider }

type collider struct {
	Layers []ColliderLayer
}

func NewCollider(layers ...ColliderLayer) Collider {
	return Collider{collider{Layers: layers}}
}
