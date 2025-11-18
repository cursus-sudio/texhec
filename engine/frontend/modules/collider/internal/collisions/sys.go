package collisions

import (
	"frontend/modules/collider"
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"
)

func NewColliderSystem(
	logger logger.Logger,
	serviceFactory ecs.ToolFactory[CollisionService],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.Query().Require(
			ecs.GetComponentType(transform.TransformComponent{}),
			ecs.GetComponentType(collider.ColliderComponent{}),
		).Build()
		service := serviceFactory.Build(w)
		query.OnAdd(func(ei []ecs.EntityID) { service.Add(ei...) })
		query.OnChange(func(ei []ecs.EntityID) { service.Update(ei...) })
		query.OnRemove(func(ei []ecs.EntityID) { service.Remove(ei...) })
		return nil
	})
}
