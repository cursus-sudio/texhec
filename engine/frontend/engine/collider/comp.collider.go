package collider

import (
	"frontend/services/assets"
)

type Collider struct{ ID assets.AssetID }

func NewCollider(id assets.AssetID) Collider {
	return Collider{ID: id}
}
