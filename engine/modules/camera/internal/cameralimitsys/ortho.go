package cameralimitsys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

type orthoSys struct {
	world     ecs.World
	transform transform.Service
	camera    camera.Service

	dirtySet ecs.DirtySet

	logger logger.Logger
}

func NewOrthoSys(
	world ecs.World,
	transform transform.Service,
	camera camera.Service,
	logger logger.Logger,
) camera.System {
	s := &orthoSys{
		world:     world,
		transform: transform,
		camera:    camera,

		dirtySet: ecs.NewDirtySet(),
		logger:   logger,
	}
	s.transform.AddDirtySet(s.dirtySet)
	s.camera.Limits().AddDirtySet(s.dirtySet)
	s.camera.Ortho().AddDirtySet(s.dirtySet)
	s.camera.Viewport().AddDirtySet(s.dirtySet)
	s.camera.NormalizedViewport().AddDirtySet(s.dirtySet)
	s.transform.AbsolutePos().BeforeGet(s.BeforeGet)

	return nil
}

func (s *orthoSys) BeforeGet() {
	ei := s.dirtySet.Get()
	if len(ei) == 0 {
		return
	}
	type save struct {
		entity ecs.EntityID
		pos    transform.AbsolutePosComponent
	}
	saves := []save{}

	for _, entity := range ei {
		limits, ok := s.camera.Limits().Get(entity)
		if !ok {
			continue
		}
		ortho, ok := s.camera.Ortho().Get(entity)
		if !ok {
			continue
		}

		pos, ok := s.transform.AbsolutePos().Get(entity)
		if !ok {
			continue
		}
		x, y, w, h := s.camera.GetViewport(entity)
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
		s.transform.AbsolutePos().Set(save.entity, save.pos)
	}
}
