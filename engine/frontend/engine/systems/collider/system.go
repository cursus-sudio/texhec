package collidersys

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/engine/tools/broadcollision"
	"shared/services/ecs"
)

func NewColliderSystem(
	serviceFactory broadcollision.CollisionServiceFactory,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.QueryEntitiesWithComponents(
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(collider.Collider{}),
		)
		query.OnAdd(func(ei []ecs.EntityID) {
			service := serviceFactory(w)
			service.Add(ei...)
		})
		query.OnChange(func(ei []ecs.EntityID) {
			service := serviceFactory(w)
			service.Update(ei...)
		})
		query.OnRemove(func(ei []ecs.EntityID) {
			service := serviceFactory(w)
			service.Remove(ei...)
		})
		return nil
	})
}
