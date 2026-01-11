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
	world     ecs.World
	transform transform.Service
	camera    camera.Service

	window window.Api
	logger logger.Logger

	perspectivesDirtySet ecs.DirtySet
	orthoDirtySet        ecs.DirtySet
}

func NewUpdateProjectionsSystem(
	eventsBuilder events.Builder,
	world ecs.World,
	transform transform.Service,
	cameraService camera.Service,
	window window.Api,
	logger logger.Logger,
) camera.System {
	return ecs.NewSystemRegister(func() error {
		s := &updateProjetionsSystem{
			world:     world,
			transform: transform,
			camera:    cameraService,

			window: window,
			logger: logger,

			perspectivesDirtySet: ecs.NewDirtySet(),
			orthoDirtySet:        ecs.NewDirtySet(),
		}

		s.camera.Perspective().BeforeGet(s.UpsertPerspective)
		s.camera.DynamicPerspective().AddDirtySet(s.perspectivesDirtySet)

		ecs.GetComponentsArray[camera.OrthoComponent](world).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.OrthoResolutionComponent](world).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.ViewportComponent](world).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.NormalizedViewportComponent](world).AddDirtySet(s.orthoDirtySet)

		s.camera.Ortho().AddDirtySet(s.orthoDirtySet)
		s.camera.Ortho().BeforeGet(s.UpsertOrtho)

		events.Listen(eventsBuilder, s.Listen)
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
		resizeOrtho, ok := s.camera.Ortho().Get(entity)
		if !ok {
			continue
		}

		x, y, w, h := s.camera.GetViewport(entity)

		size := transform.NewSize(
			float32(w-x)/resizeOrtho.Zoom,
			float32(h-y)/resizeOrtho.Zoom,
			mgl32.Abs(resizeOrtho.Far-resizeOrtho.Near),
		)
		s.transform.AbsoluteSize().Set(entity, transform.AbsoluteSizeComponent(size))
	}
}

func (s *updateProjetionsSystem) UpsertPerspective() {
	ei := s.perspectivesDirtySet.Get()
	aspectRatio := s.AspectRatio()
	for _, entity := range ei {
		resizePerspective, ok := s.camera.DynamicPerspective().Get(entity)
		if !ok {
			continue
		}
		perspective := camera.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		s.camera.Perspective().Set(entity, perspective)
	}
}

func (s *updateProjetionsSystem) Listen(e camera.ChangedResolutionEvent) {
	for _, entity := range s.camera.Ortho().GetEntities() {
		s.orthoDirtySet.Dirty(entity)
	}
	s.UpsertOrtho()

	for _, entity := range s.camera.DynamicPerspective().GetEntities() {
		s.perspectivesDirtySet.Dirty(entity)
	}
	s.UpsertPerspective()
}
