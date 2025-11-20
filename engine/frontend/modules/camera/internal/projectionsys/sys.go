package projectionsys

import (
	"frontend/modules/camera"
	"frontend/modules/transform"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

// events

// system

type updateProjetionsSystem struct {
	world  ecs.World
	window window.Api
	logger logger.Logger

	transformTool transform.TransformTool

	cameraArray              ecs.ComponentsArray[camera.CameraComponent]
	dynamicPerspectivesArray ecs.ComponentsArray[camera.DynamicPerspective]
	dynamicOrthoArray        ecs.ComponentsArray[camera.DynamicOrthoComponent]
	perspectivesArray        ecs.ComponentsArray[camera.Perspective]
	orthoArray               ecs.ComponentsArray[camera.OrthoComponent]
}

func NewUpdateProjectionsSystem(
	window window.Api,
	logger logger.Logger,
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &updateProjetionsSystem{
			world:  w,
			window: window,
			logger: logger,

			transformTool: transformToolFactory.Build(w),

			dynamicPerspectivesArray: ecs.GetComponentsArray[camera.DynamicPerspective](w.Components()),
			dynamicOrthoArray:        ecs.GetComponentsArray[camera.DynamicOrthoComponent](w.Components()),
			perspectivesArray:        ecs.GetComponentsArray[camera.Perspective](w.Components()),
			orthoArray:               ecs.GetComponentsArray[camera.OrthoComponent](w.Components()),
		}

		s.dynamicPerspectivesArray.OnAdd(s.UpsertPerspective)
		s.dynamicPerspectivesArray.OnChange(s.UpsertPerspective)

		s.dynamicOrthoArray.OnAdd(s.UpsertOrtho)
		s.dynamicOrthoArray.OnChange(s.UpsertOrtho)

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
	orthoTransaction := s.orthoArray.Transaction()
	w, h := s.Resolution()
	for _, entity := range ei {
		transform := transformTransaction.GetEntity(entity)
		size, err := transform.AbsoluteSize().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		resizeOrtho, err := s.dynamicOrthoArray.GetComponent(entity)
		if err != nil {
			continue
		}

		size.Size = mgl32.Vec3{
			w / resizeOrtho.Zoom, h / resizeOrtho.Zoom, size.Size.Z(),
		}
		transform.AbsoluteSize().Set(size)

		ortho := camera.NewOrtho(
			w, h,
			resizeOrtho.Near, resizeOrtho.Far,
			resizeOrtho.Zoom,
		)
		orthoTransaction.SaveComponent(entity, ortho)
	}

	s.logger.Warn(ecs.FlushMany(append(transformTransaction.Transactions(), orthoTransaction)...))
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
	s.UpsertOrtho(s.dynamicOrthoArray.GetEntities())
	s.UpsertPerspective(s.dynamicPerspectivesArray.GetEntities())
}
