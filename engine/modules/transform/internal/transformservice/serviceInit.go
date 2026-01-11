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
		t.absolutePosArray,
		t.absoluteRotationArray,
		t.absoluteSizeArray,
	}

	t.posArray.SetEmpty(transform.PosComponent{Pos: mgl32.Vec3{0, 0, 0}})
	t.rotationArray.SetEmpty(t.defaultRot)
	t.sizeArray.SetEmpty(t.defaultSize)

	t.maxSizeArray.SetEmpty(transform.NewMaxSize(0, 0, 0)) // 0 means not set
	t.minSizeArray.SetEmpty(transform.NewMinSize(0, 0, 0)) // 0 means not set

	t.aspectRatioArray.SetEmpty(transform.NewAspectRatio(0, 0, 0, 0)) // 0 means not set
	t.pivotPointArray.SetEmpty(t.defaultPivot)

	t.parentMaskArray.SetEmpty(transform.NewParent(transform.RelativePos))
	t.parentPivotPointArray.SetEmpty(t.defaultParentPivot)

	t.absolutePosArray.SetEmpty(transform.AbsolutePosComponent{Pos: mgl32.Vec3{0, 0, 0}})
	t.absoluteRotationArray.SetEmpty(transform.AbsoluteRotationComponent(t.defaultRot))
	t.absoluteSizeArray.SetEmpty(transform.AbsoluteSizeComponent(t.defaultSize))

	for _, arr := range arrays {
		arr.AddDependency(t.posArray)
		arr.AddDependency(t.rotationArray)
		arr.AddDependency(t.sizeArray)

		arr.AddDependency(t.maxSizeArray)
		arr.AddDependency(t.minSizeArray)

		arr.AddDependency(t.aspectRatioArray)
		arr.AddDependency(t.pivotPointArray)

		arr.AddDependency(t.hierarchy.Component())
		arr.AddDependency(t.parentMaskArray)
		arr.AddDependency(t.parentPivotPointArray)
	}

	for _, array := range arrays {
		array.AddDirtySet(t.dirtySet)
		array.BeforeGet(t.BeforeGet)
	}
}
