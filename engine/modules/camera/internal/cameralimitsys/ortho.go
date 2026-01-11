package cameralimitsys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type orthoSys struct {
	World     ecs.World         `inject:"1"`
	Transform transform.Service `inject:"1"`
	Camera    camera.Service    `inject:"1"`

	dirtySet ecs.DirtySet

	Logger logger.Logger `inject:"1"`
}

func NewOrthoSys(c ioc.Dic) camera.System {
	s := ioc.GetServices[*orthoSys](c)
	s.dirtySet = ecs.NewDirtySet()

	s.Transform.AddDirtySet(s.dirtySet)
	s.Camera.Limits().AddDirtySet(s.dirtySet)
	s.Camera.Ortho().AddDirtySet(s.dirtySet)
	s.Camera.Viewport().AddDirtySet(s.dirtySet)
	s.Camera.NormalizedViewport().AddDirtySet(s.dirtySet)
	s.Transform.AbsolutePos().BeforeGet(s.BeforeGet)

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
		limits, ok := s.Camera.Limits().Get(entity)
		if !ok {
			continue
		}
		ortho, ok := s.Camera.Ortho().Get(entity)
		if !ok {
			continue
		}

		pos, ok := s.Transform.AbsolutePos().Get(entity)
		if !ok {
			continue
		}
		x, y, w, h := s.Camera.GetViewport(entity)
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
		s.Transform.AbsolutePos().Set(save.entity, save.pos)
	}
}
