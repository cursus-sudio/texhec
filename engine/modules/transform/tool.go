package transform

import (
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type Tool interface {
	Transaction() Transaction
	Query(ecs.LiveQueryBuilder) ecs.LiveQueryBuilder
}

type Transaction interface {
	GetObject(ecs.EntityID) Object
	Transactions() []ecs.AnyComponentsArrayTransaction
	Flush() error
}

// absolute components return errors only in case when
// entity is relative to parent that doesn't exist
type Object interface {
	Pos() ecs.EntityComponent[PosComponent]
	AbsolutePos() ecs.EntityComponent[PosComponent]

	Rotation() ecs.EntityComponent[RotationComponent]
	AbsoluteRotation() ecs.EntityComponent[RotationComponent]

	Size() ecs.EntityComponent[SizeComponent]
	AbsoluteSize() ecs.EntityComponent[SizeComponent]

	MaxSize() ecs.EntityComponent[MaxSizeComponent]
	MinSize() ecs.EntityComponent[MinSizeComponent]

	AspectRatio() ecs.EntityComponent[AspectRatioComponent]
	PivotPoint() ecs.EntityComponent[PivotPointComponent]

	Parent() ecs.EntityComponent[ParentComponent]
	ParentPivotPoint() ecs.EntityComponent[ParentPivotPointComponent]

	Mat4() mgl32.Mat4
}
