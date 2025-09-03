package collidersystem

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/engine/tools/broadcollision"
	"frontend/services/ecs"
)

type ColliderSystem struct {
	world ecs.World
}

func NewColliderSystem(
	world ecs.World,
	serviceFactory broadcollision.CollisionServiceFactory,
) ColliderSystem {
	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(transform.Transform{}),
		ecs.GetComponentType(collider.Collider{}),
	)
	query.OnAdd(func(ei []ecs.EntityID) {
		service := serviceFactory(world)
		service.Add(ei...)
	})
	query.OnChange(func(ei []ecs.EntityID) {
		service := serviceFactory(world)
		service.Update(ei...)
	})
	query.OnRemove(func(ei []ecs.EntityID) {
		service := serviceFactory(world)
		service.Remove(ei...)
	})
	return ColliderSystem{world}
}
