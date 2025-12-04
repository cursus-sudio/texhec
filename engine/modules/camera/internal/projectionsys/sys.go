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

	transformTool transform.Tool
	cameraTool    camera.Tool

	cameraArray              ecs.ComponentsArray[camera.CameraComponent]
	dynamicPerspectivesArray ecs.ComponentsArray[camera.DynamicPerspective]
	perspectivesArray        ecs.ComponentsArray[camera.Perspective]
	orthoArray               ecs.ComponentsArray[camera.OrthoComponent]
}

func NewUpdateProjectionsSystem(
	window window.Api,
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.Tool],
	cameraToolFactory ecs.ToolFactory[camera.Tool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &updateProjetionsSystem{
			world:  w,
			window: window,
			logger: logger,

			transformTool: transformToolFactory.Build(w),
			cameraTool:    cameraToolFactory.Build(w),

			dynamicPerspectivesArray: ecs.GetComponentsArray[camera.DynamicPerspective](w),
			perspectivesArray:        ecs.GetComponentsArray[camera.Perspective](w),
			orthoArray:               ecs.GetComponentsArray[camera.OrthoComponent](w),
		}

		s.dynamicPerspectivesArray.OnAdd(s.UpsertPerspective)
		s.dynamicPerspectivesArray.OnChange(s.UpsertPerspective)

		orthoQuery := w.Query().
			Require(camera.OrthoComponent{}).
			Track(camera.OrthoResolutionComponent{}).
			Track(camera.ViewportComponent{}).
			Track(camera.NormalizedViewportComponent{}).
			Build()
		orthoQuery.OnAdd(s.UpsertOrtho)
		orthoQuery.OnChange(s.UpsertOrtho)

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

func (s *updateProjetionsSystem) UpsertOrtho(ei []ecs.EntityID) {
	transformTransaction := s.transformTool.Transaction()
	for _, entity := range ei {
		camera, err := s.cameraTool.GetObject(entity)
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		transformObj := transformTransaction.GetObject(entity)
		s.cameraTool.GetObject(entity)
		resizeOrtho, err := s.orthoArray.GetComponent(entity)
		if err != nil {
			continue
		}

		x, y, w, h := camera.Viewport()

		size := transform.NewSize(
			float32(w-x)/resizeOrtho.Zoom,
			float32(h-y)/resizeOrtho.Zoom,
			mgl32.Abs(resizeOrtho.Far-resizeOrtho.Near),
		)
		transformObj.AbsoluteSize().Set(size)
	}

	s.logger.Warn(transformTransaction.Flush())
}

func (s *updateProjetionsSystem) UpsertPerspective(ei []ecs.EntityID) {
	perspectiveTransaction := s.perspectivesArray.Transaction()
	aspectRatio := s.AspectRatio()
	for _, entity := range s.dynamicPerspectivesArray.GetEntities() {
		resizePerspective, err := s.dynamicPerspectivesArray.GetComponent(entity)
		if err != nil {
			continue
		}
		perspective := camera.NewPerspective(
			resizePerspective.FovY, aspectRatio,
			resizePerspective.Near, resizePerspective.Far,
		)
		perspectiveTransaction.SaveComponent(entity, perspective)
	}
	s.logger.Warn(ecs.FlushMany(perspectiveTransaction))
}

func (s *updateProjetionsSystem) Listen(e camera.ChangedResolutionEvent) {
	s.UpsertOrtho(s.orthoArray.GetEntities())
	s.UpsertPerspective(s.dynamicPerspectivesArray.GetEntities())
}
