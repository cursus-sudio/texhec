package changetransform

import (
	_ "embed"
	"frontend/engine/components/transform"
	"frontend/services/frames"
	"shared/services/ecs"
	"shared/services/logger"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type Component struct {
	T time.Duration
}

type system struct {
	World                ecs.World
	ChangeTransformArray ecs.ComponentsArray[Component]
	TransformArray       ecs.ComponentsArray[transform.Transform]
	Logger               logger.Logger
	LiveQuery            ecs.LiveQuery
}

func NewSystem(
	world ecs.World,
	logger logger.Logger,
) ecs.SystemRegister {
	liveQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(Component{}),
		ecs.GetComponentType(transform.Transform{}),
	)
	return &system{
		World:                world,
		ChangeTransformArray: ecs.GetComponentsArray[Component](world.Components()),
		TransformArray:       ecs.GetComponentsArray[transform.Transform](world.Components()),
		Logger:               logger,
		LiveQuery:            liveQuery,
	}
}

func (s *system) Register(b events.Builder) {
	events.Listen(b, s.Listen)
}

func (s *system) Listen(args frames.FrameEvent) {
	transformTransaction := s.TransformArray.Transaction()
	changeTransaction := s.ChangeTransformArray.Transaction()
	for _, entity := range s.LiveQuery.Entities() {
		changeTransformOverTimeComponent, err := s.ChangeTransformArray.GetComponent(entity)
		if err != nil {
			continue
		}
		transformComponent, err := s.TransformArray.GetComponent(entity)
		if err != nil {
			continue
		}
		changeTransformOverTimeComponent.T += args.Delta
		t := changeTransformOverTimeComponent.T

		radians := mgl32.DegToRad(float32(t.Seconds()) * 100)
		rotation := mgl32.QuatIdent().
			Mul(mgl32.QuatRotate(radians, mgl32.Vec3{1, 0, 0})).
			Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 1, 0})).
			Mul(mgl32.QuatRotate(radians, mgl32.Vec3{0, 0, 1}))
		transformComponent.Rotation = rotation

		scaleFactor := (1 + float32(t.Seconds())) / (1 + float32(t.Seconds()-args.Delta.Seconds()))
		transformComponent.Size = transformComponent.Size.Mul(scaleFactor)

		transformTransaction.SaveComponent(entity, transformComponent)
		changeTransaction.SaveComponent(entity, changeTransformOverTimeComponent)
	}
	if err := ecs.FlushMany(transformTransaction, changeTransaction); err != nil {
		s.Logger.Error(err)
	}
}
