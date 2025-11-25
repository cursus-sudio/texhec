package cameralimitsys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

type orthoSys struct {
	world         ecs.World
	query         ecs.LiveQuery
	limitsArray   ecs.ComponentsArray[camera.CameraLimitsComponent]
	orthoArray    ecs.ComponentsArray[camera.OrthoComponent]
	transformTool transform.Tool
	cameraTool    camera.Tool

	logger logger.Logger
}

func NewOrthoSys(
	transformToolFactory ecs.ToolFactory[transform.Tool],
	cameraToolFactory ecs.ToolFactory[camera.Tool],
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transformTool := transformToolFactory.Build(w)
		cameraTool := cameraToolFactory.Build(w)
		s := &orthoSys{
			world: w,
			query: transformTool.Query(w.Query()).
				Require(ecs.GetComponentType(camera.CameraLimitsComponent{})).
				Require(ecs.GetComponentType(camera.OrthoComponent{})).
				Track(ecs.GetComponentType(camera.ViewportComponent{})).
				Track(ecs.GetComponentType(camera.NormalizedViewportComponent{})).
				Build(),
			limitsArray:   ecs.GetComponentsArray[camera.CameraLimitsComponent](w),
			orthoArray:    ecs.GetComponentsArray[camera.OrthoComponent](w),
			transformTool: transformTool,
			cameraTool:    cameraTool,
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
		camera, err := s.cameraTool.GetObject(entity)
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		limits, err := s.limitsArray.GetComponent(entity)
		if err != nil {
			continue
		}
		ortho, err := s.orthoArray.GetComponent(entity)
		if err != nil {
			s.logger.Warn(err)
			continue
		}

		transform := transformTransaction.GetObject(entity)
		pos, err := transform.AbsolutePos().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		x, y, w, h := camera.Viewport()
		var halfWidth float32 = float32(w-x) / 2 / ortho.Zoom
		var halfHeight float32 = float32(h-y) / 2 / ortho.Zoom

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

		for i := 0; i < 3; i++ {
			if minPos[i] > maxPos[i] {
				center := (minPos[i] + maxPos[i]) / 2
				minPos[i], maxPos[i] = center, center
			}
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
