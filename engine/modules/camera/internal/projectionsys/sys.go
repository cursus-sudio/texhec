package projectionsys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

// events

// system

type updateProjetionsSystem struct {
	World     ecs.World         `inject:"1"`
	Transform transform.Service `inject:"1"`
	Camera    camera.Service    `inject:"1"`

	EventsBuilder events.Builder `inject:"1"`
	Window        window.Api     `inject:"1"`
	Logger        logger.Logger  `inject:"1"`

	perspectivesDirtySet ecs.DirtySet
	orthoDirtySet        ecs.DirtySet
}

func NewUpdateProjectionsSystem(c ioc.Dic) camera.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*updateProjetionsSystem](c)
		s.perspectivesDirtySet = ecs.NewDirtySet()
		s.orthoDirtySet = ecs.NewDirtySet()

		s.Camera.Perspective().BeforeGet(s.UpsertPerspective)
		s.Camera.DynamicPerspective().AddDirtySet(s.perspectivesDirtySet)

		ecs.GetComponentsArray[camera.OrthoComponent](s.World).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.OrthoResolutionComponent](s.World).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.ViewportComponent](s.World).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.NormalizedViewportComponent](s.World).AddDirtySet(s.orthoDirtySet)

		s.Camera.Ortho().AddDirtySet(s.orthoDirtySet)
		s.Camera.Ortho().BeforeGet(s.UpsertOrtho)

		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *updateProjetionsSystem) Resolution() (w float32, h float32) {
	width, height := s.Window.Window().GetSize()
	return float32(width), float32(height)
}

func (s *updateProjetionsSystem) AspectRatio() float32 {
	w, h := s.Resolution()
	return w / h
}

func (s *updateProjetionsSystem) UpsertOrtho() {
	ei := s.orthoDirtySet.Get()
	for _, entity := range ei {
		resizeOrtho, ok := s.Camera.Ortho().Get(entity)
		if !ok {
			continue
		}

		x, y, w, h := s.Camera.GetViewport(entity)

		size := transform.NewSize(
			float32(w-x)/resizeOrtho.Zoom,
			float32(h-y)/resizeOrtho.Zoom,
			mgl32.Abs(resizeOrtho.Far-resizeOrtho.Near),
		)
		s.Transform.AbsoluteSize().Set(entity, transform.AbsoluteSizeComponent(size))
	}
}

func (s *updateProjetionsSystem) UpsertPerspective() {
	ei := s.perspectivesDirtySet.Get()
	aspectRatio := s.AspectRatio()
	for _, entity := range ei {
		resizePerspective, ok := s.Camera.DynamicPerspective().Get(entity)
		if !ok {
			continue
		}
		perspective := camera.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		s.Camera.Perspective().Set(entity, perspective)
	}
}

func (s *updateProjetionsSystem) Listen(e camera.ChangedResolutionEvent) {
	for _, entity := range s.Camera.Ortho().GetEntities() {
		s.orthoDirtySet.Dirty(entity)
	}
	s.UpsertOrtho()

	for _, entity := range s.Camera.DynamicPerspective().GetEntities() {
		s.perspectivesDirtySet.Dirty(entity)
	}
	s.UpsertPerspective()
}
