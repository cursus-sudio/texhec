package collisions

import (
	"engine/modules/collider"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
)

func NewColliderSystem(
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.Tool],
	serviceFactory ecs.ToolFactory[CollisionService],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transformTool := transformToolFactory.Build(w)
		query := transformTool.Query(w.Query()).
			Require(collider.ColliderComponent{}).
			Build()
		service := serviceFactory.Build(w)
		query.OnAdd(func(ei []ecs.EntityID) { service.Add(ei...) })
		query.OnChange(func(ei []ecs.EntityID) { service.Update(ei...) })
		query.OnRemove(func(ei []ecs.EntityID) { service.Remove(ei...) })
		return nil
	})
}
