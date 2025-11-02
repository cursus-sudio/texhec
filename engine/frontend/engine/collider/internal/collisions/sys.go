package collisions

import (
	"frontend/engine/collider"
	"frontend/engine/transform"
	"shared/services/ecs"
)

func NewColliderSystem(
	serviceFactory ecs.ToolFactory[CollisionService],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.Query().Require(
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(collider.Collider{}),
		).Build()
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
