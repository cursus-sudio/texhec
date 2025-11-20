package cameralimitsys

import (
	"frontend/modules/camera"
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

type orthoSys struct {
	world         ecs.World
	query         ecs.LiveQuery
	limitsArray   ecs.ComponentsArray[camera.CameraLimitsComponent]
	orthoArray    ecs.ComponentsArray[camera.OrthoComponent]
	transformTool transform.TransformTool

	logger logger.Logger
}

func NewOrthoSys(
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transformTool := transformToolFactory.Build(w)
		s := &orthoSys{
			world: w,
			query: transformTool.Query(w.Query()).
				Require(ecs.GetComponentType(camera.CameraLimitsComponent{})).
				Require(ecs.GetComponentType(camera.OrthoComponent{})).
				Build(),
			limitsArray:   ecs.GetComponentsArray[camera.CameraLimitsComponent](w.Components()),
			orthoArray:    ecs.GetComponentsArray[camera.OrthoComponent](w.Components()),
			transformTool: transformTool,
			logger:        logger,
		}
		s.Addlisteners()
		return nil
	})
}

func (s *orthoSys) Addlisteners() {
	s.query.OnAdd(s.ChangeListener)
	s.query.OnChange(s.ChangeListener)
}

func (s *orthoSys) ChangeListener(ei []ecs.EntityID) {
	transformTransaction := s.transformTool.Transaction()
	for _, entity := range ei {
		limits, err := s.limitsArray.GetComponent(entity)
		if err != nil {
			continue
		}
		orthoComponent, err := s.orthoArray.GetComponent(entity)
		if err != nil {
			continue
		}

		transform := transformTransaction.GetEntity(entity)
		pos, err := transform.AbsolutePos().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
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

		pos.Pos = mgl32.Vec3{
			max(pos.Pos.X(), minPos.X()),
			max(pos.Pos.Y(), minPos.Y()),
			max(pos.Pos.Z(), minPos.Z()),
		}

		pos.Pos = mgl32.Vec3{
			min(pos.Pos.X(), maxPos.X()),
			min(pos.Pos.Y(), maxPos.Y()),
			min(pos.Pos.Z(), maxPos.Z()),
		}

		transform.AbsolutePos().Set(pos)
	}

	s.logger.Warn(ecs.FlushMany(transformTransaction.Transactions()...))
}
