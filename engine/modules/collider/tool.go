package collider

import (
	"engine/modules/groups"
	"engine/modules/transform"
	"engine/services/ecs"
)

type ColliderTool interface {
	Collider() Interface
}

type World interface {
	ecs.World
	transform.TransformTool
	groups.GroupsTool
}

type Interface interface {
	Component() ecs.ComponentsArray[Component]

	// todo add collision groups
	// narrow
	CollidesWithRay(ecs.EntityID, Ray) (ObjectRayCollision, error)
	CollidesWithObject(ecs.EntityID, ecs.EntityID) (ObjectObjectCollision, error)

	// broad
	ShootRay(Ray) (ObjectRayCollision, error)
	NarrowCollisions(ecs.EntityID) ([]ecs.EntityID, error)
}
