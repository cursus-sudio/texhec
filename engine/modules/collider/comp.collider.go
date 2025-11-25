package collider

import (
	"engine/services/assets"
)

type ColliderComponent struct{ ID assets.AssetID }

func NewCollider(id assets.AssetID) ColliderComponent {
	return ColliderComponent{ID: id}
}
