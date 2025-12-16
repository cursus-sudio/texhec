package transformtool

import (
	"engine/modules/transform"
	"engine/services/ecs"
)

type save struct {
	entity ecs.EntityID
	pos    transform.AbsolutePosComponent
	rot    transform.AbsoluteRotationComponent
	size   transform.AbsoluteSizeComponent
}

func (t tool) Init() {
	arrays := []ecs.AnyComponentArray{
		t.absolutePosArray,
		t.absoluteRotationArray,
		t.absoluteSizeArray,
	}

	for _, arr := range arrays {
		arr.AddDependency(t.posArray)
		arr.AddDependency(t.rotationArray)
		arr.AddDependency(t.sizeArray)

		arr.AddDependency(t.maxSizeArray)
		arr.AddDependency(t.minSizeArray)

		arr.AddDependency(t.aspectRatioArray)
		arr.AddDependency(t.pivotPointArray)

		arr.AddDependency(t.hierarchyArray)
		arr.AddDependency(t.parentMaskArray)
		arr.AddDependency(t.parentPivotPointArray)
	}

	for _, array := range arrays {
		array.AddDirtySet(t.dirtySet)
		array.BeforeGet(t.BeforeGet)
	}
}
