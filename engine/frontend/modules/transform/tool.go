package transform

import (
	"shared/services/datastructures"
	"shared/services/ecs"

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

	PivotPoint() ecs.EntityComponent[PivotPointComponent]

	Parent() ecs.EntityComponent[ParentComponent]
	ParentPivotPoint() ecs.EntityComponent[ParentPivotPointComponent]

	Children() datastructures.SparseSetReader[ecs.EntityID]
	// includes children of children
	FlatChildren() datastructures.SparseSetReader[ecs.EntityID]

	Mat4() mgl32.Mat4
}
