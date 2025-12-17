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
	world  ecs.World
	window window.Api
	logger logger.Logger

	transformTool transform.Interface
	cameraTool    camera.Interface

	perspectivesDirtySet ecs.DirtySet
	orthoDirtySet        ecs.DirtySet

	cameraArray              ecs.ComponentsArray[camera.CameraComponent]
	dynamicPerspectivesArray ecs.ComponentsArray[camera.DynamicPerspective]
	perspectivesArray        ecs.ComponentsArray[camera.Perspective]
	orthoArray               ecs.ComponentsArray[camera.OrthoComponent]
}

func NewUpdateProjectionsSystem(
	window window.Api,
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	cameraToolFactory ecs.ToolFactory[camera.CameraTool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &updateProjetionsSystem{
			world:  w,
			window: window,
			logger: logger,

			transformTool: transformToolFactory.Build(w).Transform(),
			cameraTool:    cameraToolFactory.Build(w).Camera(),

			perspectivesDirtySet: ecs.NewDirtySet(),
			orthoDirtySet:        ecs.NewDirtySet(),

			dynamicPerspectivesArray: ecs.GetComponentsArray[camera.DynamicPerspective](w),
			perspectivesArray:        ecs.GetComponentsArray[camera.Perspective](w),
			orthoArray:               ecs.GetComponentsArray[camera.OrthoComponent](w),
		}

		s.perspectivesArray.BeforeGet(s.UpsertPerspective)
		s.dynamicPerspectivesArray.AddDirtySet(s.perspectivesDirtySet)

		ecs.GetComponentsArray[camera.OrthoComponent](w).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.OrthoResolutionComponent](w).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.ViewportComponent](w).AddDirtySet(s.orthoDirtySet)
		ecs.GetComponentsArray[camera.NormalizedViewportComponent](w).AddDirtySet(s.orthoDirtySet)

		s.orthoArray.AddDirtySet(s.orthoDirtySet)
		s.orthoArray.BeforeGet(s.UpsertOrtho)

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
		camera, err := s.cameraTool.GetObject(entity)
		if err != nil {
			continue
		}
		s.cameraTool.GetObject(entity)
		resizeOrtho, ok := s.orthoArray.Get(entity)
		if !ok {
			continue
		}

		x, y, w, h := camera.Viewport()

		size := transform.NewSize(
			float32(w-x)/resizeOrtho.Zoom,
			float32(h-y)/resizeOrtho.Zoom,
			mgl32.Abs(resizeOrtho.Far-resizeOrtho.Near),
		)
		s.transformTool.SetAbsoluteSize(entity, transform.AbsoluteSizeComponent(size))
	}
}

func (s *updateProjetionsSystem) UpsertPerspective() {
	ei := s.perspectivesDirtySet.Get()
	aspectRatio := s.AspectRatio()
	for _, entity := range ei {
		resizePerspective, ok := s.dynamicPerspectivesArray.Get(entity)
		if !ok {
			continue
		}
		perspective := camera.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		s.perspectivesArray.Set(entity, perspective)
	}
}

func (s *updateProjetionsSystem) Listen(e camera.ChangedResolutionEvent) {
	for _, entity := range s.orthoArray.GetEntities() {
		s.orthoDirtySet.Dirty(entity)
	}
	s.UpsertOrtho()

	for _, entity := range s.dynamicPerspectivesArray.GetEntities() {
		s.perspectivesDirtySet.Dirty(entity)
	}
	s.UpsertPerspective()
}
