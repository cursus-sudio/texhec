package collider

import (
	"engine/modules/assets"
)

type Component struct{ ID assets.ID }

func NewCollider(id assets.ID) Component {
	return Component{ID: id}
}
