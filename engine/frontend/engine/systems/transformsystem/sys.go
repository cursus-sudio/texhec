package transformsystem

import (
	"frontend/engine/components/transform"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

func NewPivotPointSystem(world ecs.World, logger logger.Logger) {
	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(transform.Transform{}),
		ecs.GetComponentType(transform.PivotPoint{}),
	)
	transformArray := ecs.GetComponentsArray[transform.Transform](world.Components())
	transformTransaction := transformArray.Transaction()
	pivotPointsArray := ecs.GetComponentsArray[transform.PivotPoint](world.Components())
	listener := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			transformComponent, err := transformArray.GetComponent(entity)
			if err != nil {
				continue
			}
			pivot, err := pivotPointsArray.GetComponent(entity)
			if err != nil {
				continue
			}

			change := mgl32.Vec3{
				transformComponent.Size[0] * (pivot.Point[0] - .5),
				transformComponent.Size[1] * (pivot.Point[1] - .5),
				transformComponent.Size[2] * (pivot.Point[2] - .5),
			}
			transformComponent.SetPos(transformComponent.Pos.Add(change))
			transformTransaction.DirtySaveComponent(entity, transformComponent)
		}
		if err := transformTransaction.Flush(); err != nil {
			logger.Error(err)
		}
	}
	query.OnAdd(listener)
	query.OnChange(listener)
}
