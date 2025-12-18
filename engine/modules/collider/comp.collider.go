package collider

import (
	"engine/services/assets"
)

type Component struct{ ID assets.AssetID }

func NewCollider(id assets.AssetID) Component {
	return Component{ID: id}
}
