package internal

import (
	"frontend/modules/camera"
	"frontend/modules/drag"
	"frontend/modules/transform"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type s struct {
	camera ecs.ToolFactory[camera.CameraTool]
}

func NewSystem(
	camera ecs.ToolFactory[camera.CameraTool],
) ecs.SystemRegister {
	return s{camera}
}

func (s s) Register(w ecs.World) error {
	transformArray := ecs.GetComponentsArray[transform.TransformComponent](w.Components())
	cameraTool := s.camera.Build(w)
	events.Listen(w.EventsBuilder(), func(event drag.DraggableEvent) {
		camera, err := cameraTool.Get(event.Drag.Camera)
		if err != nil {
			return
		}

		transform, err := transformArray.GetComponent(event.Entity)
		if err != nil {
			return
		}

		fromRay := camera.ShootRay(event.Drag.From)
		toRay := camera.ShootRay(event.Drag.To)

		posDiff := toRay.Pos.Sub(fromRay.Pos)
		transform.Pos = transform.Pos.Add(posDiff)

		rotDiff := mgl32.QuatBetweenVectors(toRay.Direction, fromRay.Direction)
		transform.Rotation = transform.Rotation.Mul(rotDiff)
		transformArray.SaveComponent(event.Entity, transform)
	})
	return nil
}
