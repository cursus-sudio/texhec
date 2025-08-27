package broadcollision

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders/shapes"
	"frontend/services/ecs"
)

type BroadCollisionFactory interface {
	Generate(world ecs.World) BroadCollision
}

type BroadCollider struct {
	ID       ecs.EntityID
	Collider transform.AABB
}

type BroadCollision interface {
	Upsert(entities []BroadCollider)
	Remove(entities []ecs.EntityID)
	GetColliding(aabb transform.AABB) []ecs.EntityID
	GetRayColliding(ray shapes.Ray) []ecs.EntityID
}
