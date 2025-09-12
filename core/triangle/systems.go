package triangle

import (
	_ "embed"
	"frontend/engine/components/transform"
	"frontend/services/ecs"
	"frontend/services/frames"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

type ChangeTransformOverTimeComponent struct {
	T time.Duration
}

type ChangeTransformOverTimeSystem struct {
	World                ecs.World
	ChangeTransformArray ecs.ComponentsArray[ChangeTransformOverTimeComponent]
	TransformArray       ecs.ComponentsArray[transform.Transform]
	LiveQuery            ecs.LiveQuery
}

func NewChangeTransformOverTimeSystem(
	world ecs.World,
) *ChangeTransformOverTimeSystem {
	liveQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(ChangeTransformOverTimeComponent{}),
		ecs.GetComponentType(transform.Transform{}),
	)
	return &ChangeTransformOverTimeSystem{
		World:                world,
		ChangeTransformArray: ecs.GetComponentsArray[ChangeTransformOverTimeComponent](world.Components()),
		TransformArray:       ecs.GetComponentsArray[transform.Transform](world.Components()),
		LiveQuery:            liveQuery,
	}
}

func (s *ChangeTransformOverTimeSystem) Update(args frames.FrameEvent) {
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
		// transformComponent.Size[0]*= scaleFactor
		// transformComponent.Size.Y *= scaleFactor
		// transformComponent.Size.Z *= scaleFactor
		// transformComponent.Pos.X = float32(t.Seconds()) * 100

		s.TransformArray.SaveComponent(entity, transformComponent)
		s.ChangeTransformArray.SaveComponent(entity, changeTransformOverTimeComponent)
	}
}
