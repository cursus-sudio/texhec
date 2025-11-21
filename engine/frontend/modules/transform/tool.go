package transform

import (
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type TransformTool interface {
	Transaction() TransformTransaction
	Query(ecs.LiveQueryBuilder) ecs.LiveQueryBuilder
}

type TransformTransaction interface {
	GetEntity(ecs.EntityID) EntityTransform
	Transactions() []ecs.AnyComponentsArrayTransaction
	Flush() error
}

// absolute components return errors only in case when
// entity is relative to parent that doesn't exist
type EntityTransform interface {
	Pos() ecs.EntityComponent[PosComponent]
	AbsolutePos() ecs.EntityComponent[PosComponent]

	Rotation() ecs.EntityComponent[RotationComponent]
	AbsoluteRotation() ecs.EntityComponent[RotationComponent]

	Size() ecs.EntityComponent[SizeComponent]
	AbsoluteSize() ecs.EntityComponent[SizeComponent]

	PivotPoint() ecs.EntityComponent[PivotPointComponent]

	Parent() ecs.EntityComponent[ParentComponent]
	ParentPivotPoint() ecs.EntityComponent[ParentPivotPointComponent]

	Mat4() mgl32.Mat4
}
