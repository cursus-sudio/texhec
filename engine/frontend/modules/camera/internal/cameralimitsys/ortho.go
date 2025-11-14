package cameralimitsys

import (
	"frontend/modules/camera"
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

type orthoSys struct {
	world          ecs.World
	query          ecs.LiveQuery
	limitsArray    ecs.ComponentsArray[camera.CameraLimitsComponent]
	orthoArray     ecs.ComponentsArray[camera.OrthoComponent]
	transformArray ecs.ComponentsArray[transform.TransformComponent]

	logger logger.Logger
}

func NewOrthoSys(w ecs.World, logger logger.Logger) {
	s := &orthoSys{
		world: w,
		query: w.Query().
			Require(ecs.GetComponentType(camera.CameraLimitsComponent{})).
			Require(ecs.GetComponentType(camera.OrthoComponent{})).
			Track(ecs.GetComponentType(transform.TransformComponent{})).
			Build(),
		limitsArray:    ecs.GetComponentsArray[camera.CameraLimitsComponent](w.Components()),
		orthoArray:     ecs.GetComponentsArray[camera.OrthoComponent](w.Components()),
		transformArray: ecs.GetComponentsArray[transform.TransformComponent](w.Components()),
		logger:         logger,
	}
	s.Addlisteners()
}

func (s *orthoSys) Addlisteners() {
	s.query.OnAdd(s.ChangeListener)
	s.query.OnChange(s.ChangeListener)
}

func (s *orthoSys) ChangeListener(ei []ecs.EntityID) {
	transformTransaction := s.transformArray.Transaction()
	for _, entity := range ei {
		limits, err := s.limitsArray.GetComponent(entity)
		if err != nil {
			continue
		}
		orthoComponent, err := s.orthoArray.GetComponent(entity)
		if err != nil {
			continue
		}

		transformComponent, err := s.transformArray.GetComponent(entity)
		if err != nil {
			transformComponent = transform.NewTransform()
		}

		halfWidth := orthoComponent.Width / 2.0
		halfHeight := orthoComponent.Height / 2.0

		minPos := mgl32.Vec3{
			limits.Min.X() + halfWidth,
			limits.Min.Y() + halfHeight,
			limits.Min.Z(),
		}

		maxPos := mgl32.Vec3{
			limits.Max.X() - halfWidth,
			limits.Max.Y() - halfHeight,
			limits.Max.Z(),
		}

		transformComponent.Pos = mgl32.Vec3{
			max(transformComponent.Pos.X(), minPos.X()),
			max(transformComponent.Pos.Y(), minPos.Y()),
			max(transformComponent.Pos.Z(), minPos.Z()),
		}

		transformComponent.Pos = mgl32.Vec3{
			min(transformComponent.Pos.X(), maxPos.X()),
			min(transformComponent.Pos.Y(), maxPos.Y()),
			min(transformComponent.Pos.Z(), maxPos.Z()),
		}

		transformTransaction.DirtySaveComponent(entity, transformComponent)
	}

	s.logger.Warn(transformTransaction.Flush())
}
