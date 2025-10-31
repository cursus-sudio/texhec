package transformsys

import (
	"frontend/engine/components/transform"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

func NewPivotPointSystem(logger logger.Logger) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.Query().Require(
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(transform.PivotPoint{}),
		).Build()
		transformArray := ecs.GetComponentsArray[transform.Transform](w.Components())
		transformTransaction := transformArray.Transaction()
		pivotPointsArray := ecs.GetComponentsArray[transform.PivotPoint](w.Components())
		listener := func(ei []ecs.EntityID) {
			for _, entity := range ei {
				transformComponent, err := transformArray.GetComponent(entity)
				if err != nil {
					transformComponent = transform.NewTransform()
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
		return nil
	})
}
