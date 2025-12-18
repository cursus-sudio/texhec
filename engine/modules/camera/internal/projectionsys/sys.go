package projectionsys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

// events

// system

type updateProjetionsSystem struct {
	camera.World
	camera.CameraTool

	window window.Api
	logger logger.Logger

	perspectivesDirtySet ecs.DirtySet
	orthoDirtySet        ecs.DirtySet
}

func NewUpdateProjectionsSystem(
	window window.Api,
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.World, camera.CameraTool],
) ecs.SystemRegister[camera.World] {
	return ecs.NewSystemRegister(func(w camera.World) error {
		s := &updateProjetionsSystem{
			World:      w,
			CameraTool: cameraToolFactory.Build(w),

			window: window,
			logger: logger,

			perspectivesDirtySet: ecs.NewDirtySet(),
			orthoDirtySet:        ecs.NewDirtySet(),
		}

		s.Camera().Perspective().BeforeGet(s.UpsertPerspective)
		s.Camera().DynamicPerspective().AddDirtySet(s.perspectivesDirtySet)

		ecs.GetComponentsArray[camera.OrthoComponent](w).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.OrthoResolutionComponent](w).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.ViewportComponent](w).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.NormalizedViewportComponent](w).AddDirtySet(s.orthoDirtySet)

		s.Camera().Ortho().AddDirtySet(s.orthoDirtySet)
		s.Camera().Ortho().BeforeGet(s.UpsertOrtho)

		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *updateProjetionsSystem) Resolution() (w float32, h float32) {
	width, height := s.window.Window().GetSize()
	return float32(width), float32(height)
}

func (s *updateProjetionsSystem) AspectRatio() float32 {
	w, h := s.Resolution()
	return w / h
}

func (s *updateProjetionsSystem) UpsertOrtho() {
	ei := s.orthoDirtySet.Get()
	for _, entity := range ei {
		camera, err := s.Camera().GetObject(entity)
		if err != nil {
			continue
		}
		s.Camera().GetObject(entity)
		resizeOrtho, ok := s.Camera().Ortho().Get(entity)
		if !ok {
			continue
		}

		x, y, w, h := camera.Viewport()

		size := transform.NewSize(
			float32(w-x)/resizeOrtho.Zoom,
			float32(h-y)/resizeOrtho.Zoom,
			mgl32.Abs(resizeOrtho.Far-resizeOrtho.Near),
		)
		s.Transform().SetAbsoluteSize(entity, transform.AbsoluteSizeComponent(size))
	}
}

func (s *updateProjetionsSystem) UpsertPerspective() {
	ei := s.perspectivesDirtySet.Get()
	aspectRatio := s.AspectRatio()
	for _, entity := range ei {
		resizePerspective, ok := s.Camera().DynamicPerspective().Get(entity)
		if !ok {
			continue
		}
		perspective := camera.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		s.Camera().Perspective().Set(entity, perspective)
	}
}

func (s *updateProjetionsSystem) Listen(e camera.ChangedResolutionEvent) {
	for _, entity := range s.Camera().Ortho().GetEntities() {
		s.orthoDirtySet.Dirty(entity)
	}
	s.UpsertOrtho()

	for _, entity := range s.Camera().DynamicPerspective().GetEntities() {
		s.perspectivesDirtySet.Dirty(entity)
	}
	s.UpsertPerspective()
}
