package transformservice

import (
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type save struct {
	entity ecs.EntityID
	pos    transform.AbsolutePosComponent
	rot    transform.AbsoluteRotationComponent
	size   transform.AbsoluteSizeComponent
}

func (t *service) Init() {
	arrays := []ecs.AnyComponentArray{
		t.AbsolutePosArray,
		t.AbsoluteRotationArray,
		t.AbsoluteSizeArray,
	}

	t.PosArray.SetEmpty(transform.PosComponent{Pos: mgl32.Vec3{0, 0, 0}})
	t.RotationArray.SetEmpty(t.defaultRot)
	t.SizeArray.SetEmpty(t.defaultSize)

	t.MaxSizeArray.SetEmpty(transform.NewMaxSize(0, 0, 0)) // 0 means not set
	t.MinSizeArray.SetEmpty(transform.NewMinSize(0, 0, 0)) // 0 means not set

	t.AspectRatioArray.SetEmpty(transform.NewAspectRatio(0, 0, 0, 0)) // 0 means not set
	t.PivotPointArray.SetEmpty(t.defaultPivot)

	t.ParentMaskArray.SetEmpty(transform.NewParent(transform.RelativePos))
	t.ParentPivotPointArray.SetEmpty(t.defaultParentPivot)

	t.AbsolutePosArray.SetEmpty(transform.AbsolutePosComponent{Pos: mgl32.Vec3{0, 0, 0}})
	t.AbsoluteRotationArray.SetEmpty(transform.AbsoluteRotationComponent(t.defaultRot))
	t.AbsoluteSizeArray.SetEmpty(transform.AbsoluteSizeComponent(t.defaultSize))

	for _, arr := range arrays {
		arr.AddDependency(t.PosArray)
		arr.AddDependency(t.RotationArray)
		arr.AddDependency(t.SizeArray)

		arr.AddDependency(t.MaxSizeArray)
		arr.AddDependency(t.MinSizeArray)

		arr.AddDependency(t.AspectRatioArray)
		arr.AddDependency(t.PivotPointArray)

		arr.AddDependency(t.Hierarchy.Component())
		arr.AddDependency(t.ParentMaskArray)
		arr.AddDependency(t.ParentPivotPointArray)
	}

	for _, array := range arrays {
		array.AddDirtySet(t.DirtySet)
		array.BeforeGet(t.BeforeGet)
	}
}
