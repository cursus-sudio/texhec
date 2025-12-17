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
	dirtySet      ecs.DirtySet
	limitsArray   ecs.ComponentsArray[camera.CameraLimitsComponent]
	orthoArray    ecs.ComponentsArray[camera.OrthoComponent]
	transformTool transform.Interface
	cameraTool    camera.Interface

	logger logger.Logger
}

func NewOrthoSys(
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	cameraToolFactory ecs.ToolFactory[camera.CameraTool],
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transformTool := transformToolFactory.Build(w).Transform()
		cameraTool := cameraToolFactory.Build(w).Camera()

		dirtySet := ecs.NewDirtySet()
		transformTool.AddDirtySet(dirtySet)
		ecs.GetComponentsArray[camera.CameraLimitsComponent](w).AddDirtySet(dirtySet)
		ecs.GetComponentsArray[camera.OrthoComponent](w).AddDirtySet(dirtySet)
		ecs.GetComponentsArray[camera.ViewportComponent](w).AddDirtySet(dirtySet)
		ecs.GetComponentsArray[camera.NormalizedViewportComponent](w).AddDirtySet(dirtySet)

		s := &orthoSys{
			world:         w,
			dirtySet:      dirtySet,
			limitsArray:   ecs.GetComponentsArray[camera.CameraLimitsComponent](w),
			orthoArray:    ecs.GetComponentsArray[camera.OrthoComponent](w),
			transformTool: transformTool,
			cameraTool:    cameraTool,
			logger:        logger,
		}
		s.transformTool.Pos().BeforeGet(s.BeforeGet)

		return nil
	})
}

func (s *orthoSys) BeforeGet() {
	ei := s.dirtySet.Get()
	if len(ei) == 0 {
		return
	}
	type save struct {
		entity ecs.EntityID
		pos    transform.PosComponent
	}
	saves := []save{}

	for _, entity := range ei {
		camera, err := s.cameraTool.GetObject(entity)
		if err != nil {
			continue
		}
		limits, ok := s.limitsArray.GetComponent(entity)
		if !ok {
			continue
		}
		ortho, ok := s.orthoArray.GetComponent(entity)
		if !ok {
			continue
		}

		pos, ok := s.transformTool.Pos().GetComponent(entity)
		if !ok {
			continue
		}
		x, y, w, h := camera.Viewport()
		halfWidth := float32(w-x) / 2 / ortho.Zoom
		halfHeight := float32(h-y) / 2 / ortho.Zoom

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

		saves = append(saves, save{entity, pos})
	}
	for _, save := range saves {
		s.transformTool.Pos().SaveComponent(save.entity, save.pos)
	}
}
