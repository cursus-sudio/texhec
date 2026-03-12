package collider

import (
	"engine/services/ecs"
)

type Component struct{ ID ecs.EntityID }

func NewCollider(id ecs.EntityID) Component {
	return Component{ID: id}
}
