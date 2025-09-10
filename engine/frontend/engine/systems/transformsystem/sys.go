package transformsystem

import (
	"frontend/engine/components/transform"
	"frontend/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

func applyLockTransform(transformComponent transform.Transform, pivot transform.PivotPoint) {
	change := mgl32.Vec3{
		transformComponent.Size[0] * (pivot.Point[0] - .5),
		transformComponent.Size[1] * (pivot.Point[1] - .5),
		transformComponent.Size[2] * (pivot.Point[2] - .5),
	}
	transformComponent.Pos = transformComponent.Pos.Add(change)
}

func NewPivotPointSystem(world ecs.World) {
	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(transform.Transform{}),
		ecs.GetComponentType(transform.PivotPoint{}),
	)
	transformArray := ecs.GetComponentArray[transform.Transform](world.Components())
	pivotPointsArray := ecs.GetComponentArray[transform.PivotPoint](world.Components())
	var lastChanged ecs.EntityID
	listener := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			if entity == lastChanged {
				var zero ecs.EntityID
				lastChanged = zero
				continue
			}
			transformComponent, err := transformArray.GetComponent(entity)
			if err != nil {
				continue
			}
			pivot, err := pivotPointsArray.GetComponent(entity)
			if err != nil {
				continue
			}
			applyLockTransform(transformComponent, pivot)
		}
	}
	query.OnAdd(listener)
	query.OnChange(listener)
}
