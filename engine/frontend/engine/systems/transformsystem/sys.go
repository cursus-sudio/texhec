package transformsystem

import (
	"frontend/engine/components/transform"
	"frontend/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

func applyLockTransform(transformComponent transform.Transform, posLock transform.PosLock) {
	change := mgl32.Vec3{
		transformComponent.Size[0] * (posLock.Lock[0] - .5),
		transformComponent.Size[1] * (posLock.Lock[1] - .5),
		transformComponent.Size[2] * (posLock.Lock[2] - .5),
	}
	transformComponent.Pos = transformComponent.Pos.Add(change)
}

func NewLockerSystem(world ecs.World) {
	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(transform.Transform{}),
		ecs.GetComponentType(transform.PosLock{}),
	)
	var lastChanged ecs.EntityID
	listener := func(ei []ecs.EntityID) {
		for _, entity := range ei {
			if entity == lastChanged {
				var zero ecs.EntityID
				lastChanged = zero
				continue
			}
			transformComponent, err := ecs.GetComponent[transform.Transform](world, entity)
			if err != nil {
				continue
			}
			posLock, err := ecs.GetComponent[transform.PosLock](world, entity)
			if err != nil {
				continue
			}
			applyLockTransform(transformComponent, posLock)
		}
	}
	query.OnAdd(listener)
	query.OnChange(listener)
}
