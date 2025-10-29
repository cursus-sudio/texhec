package collidersys

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/engine/tools/broadcollision"
	"shared/services/ecs"
)

func NewColliderSystem(
	serviceFactory ecs.ToolFactory[broadcollision.CollisionService],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.QueryEntitiesWithComponents(
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(collider.Collider{}),
		)
		service := serviceFactory.Build(w)
		query.OnAdd(func(ei []ecs.EntityID) {
			service.Add(ei...)
		})
		query.OnChange(func(ei []ecs.EntityID) {
			service.Update(ei...)
		})
		query.OnRemove(func(ei []ecs.EntityID) {
			service.Remove(ei...)
		})
		return nil
	})
}
