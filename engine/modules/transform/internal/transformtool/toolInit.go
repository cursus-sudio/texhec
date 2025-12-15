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

	dirtySet := ecs.NewDirtySet()

	beforeGet := func() {
		entities := dirtySet.Get()
		if len(entities) == 0 {
			return
		}
		children := []ecs.EntityID{}

		saves := []save{}

		for len(entities) != 0 || len(children) != 0 {
			if len(entities) == 0 {
				entities = children
				for _, save := range saves {
					t.absolutePosArray.SaveComponent(save.entity, save.pos)
					t.absoluteRotationArray.SaveComponent(save.entity, save.rot)
					t.absoluteSizeArray.SaveComponent(save.entity, save.size)
				}
				dirtySet.Clear()

				children = nil
				saves = nil
			}
			entity := entities[0]
			entities = entities[1:]

			saves = append(saves, save{
				entity: entity,
				pos:    t.CalculateAbsolutePos(entity),
				rot:    t.CalculateAbsoluteRot(entity),
				size:   t.CalculateAbsoluteSize(entity),
			})

			for _, child := range t.hierarchy.Children(entity).GetIndices() {
				comparedMask := transform.RelativePos | transform.RelativeRotation | transform.RelativeSizeXYZ
				mask, ok := t.parentMaskArray.GetComponent(child)
				if !ok || mask.RelativeMask&comparedMask == 0 {
					continue
				}
				children = append(children, child)
			}
		}

		for _, save := range saves {
			t.absolutePosArray.SaveComponent(save.entity, save.pos)
			t.absoluteRotationArray.SaveComponent(save.entity, save.rot)
			t.absoluteSizeArray.SaveComponent(save.entity, save.size)
		}
		dirtySet.Clear()
	}

	for _, array := range arrays {
		array.AddDirtySet(dirtySet)
		array.BeforeGet(beforeGet)
	}
}
