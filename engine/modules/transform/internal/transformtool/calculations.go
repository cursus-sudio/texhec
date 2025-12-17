package transformtool

import (
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

func (t tool) GetRelativeParentPos(entity ecs.EntityID) mgl32.Vec3 {
	parent, ok := t.hierarchyArray.Get(entity)
	if !ok {
		return mgl32.Vec3{}
	}
	parentMask, ok := t.parentMaskArray.Get(entity)
	if !ok || parentMask.RelativeMask&transform.RelativePos == 0 {
		return mgl32.Vec3{}
	}
	parentPos, ok := t.absolutePosArray.Get(parent.Parent)
	if !ok {
		parentPos.Pos = t.CalculateAbsolutePos(parent.Parent).Pos
	}
	parentSize, ok := t.absoluteSizeArray.Get(parent.Parent)
	if !ok {
		parentSize.Size = t.CalculateAbsoluteSize(parent.Parent).Size
	}
	parentPivot, ok := t.parentPivotPointArray.Get(entity)
	if !ok {
		parentPivot = t.defaultParentPivot
	}
	parentPivot.Point = parentPivot.Point.Sub(t.defaultParentPivot.Point)
	return parentPos.Pos.Add(mgl32.Vec3{
		parentSize.Size[0] * parentPivot.Point[0],
		parentSize.Size[1] * parentPivot.Point[1],
		parentSize.Size[2] * parentPivot.Point[2],
	})
}

func (t tool) GetPivotPos(entity ecs.EntityID) mgl32.Vec3 {
	pivot, ok := t.pivotPointArray.Get(entity)
	if !ok {
		return mgl32.Vec3{}
	}
	size, ok := t.absoluteSizeArray.Get(entity)
	if !ok {
		return mgl32.Vec3{}
	}
	pivot.Point = pivot.Point.Sub(t.defaultPivot.Point)
	return mgl32.Vec3{
		size.Size[0] * (-pivot.Point[0]),
		size.Size[1] * (-pivot.Point[1]),
		size.Size[2] * (-pivot.Point[2]),
	}
}

//

func (t tool) GetRelativeParentRotation(entity ecs.EntityID) mgl32.Quat {
	parent, ok := t.hierarchyArray.Get(entity)
	if !ok {
		return mgl32.QuatIdent()
	}
	parentMask, ok := t.parentMaskArray.Get(entity)
	if !ok || parentMask.RelativeMask&transform.RelativeRotation == 0 {
		return mgl32.QuatIdent()
	}
	parentRot, ok := t.absoluteRotationArray.Get(parent.Parent)
	if !ok {
		return mgl32.QuatIdent()
	}
	return parentRot.Rotation
}

//

func (t tool) GetRelativeParentSize(entity ecs.EntityID) mgl32.Vec3 {
	size := mgl32.Vec3{1, 1, 1}
	parent, ok := t.hierarchyArray.Get(entity)
	if !ok {
		return size
	}
	parentMask, ok := t.parentMaskArray.Get(entity)
	if !ok {
		return size
	}
	parentSize, ok := t.absoluteSizeArray.Get(parent.Parent)
	if !ok {
		return size
	}
	if parentMask.RelativeMask&transform.RelativeSizeX != 0 {
		size[0] = parentSize.Size[0]
	}
	if parentMask.RelativeMask&transform.RelativeSizeY != 0 {
		size[1] = parentSize.Size[1]
	}
	if parentMask.RelativeMask&transform.RelativeSizeZ != 0 {
		size[2] = parentSize.Size[2]
	}
	return size
}

func (t tool) ApplyMinMaxSize(entity ecs.EntityID, size *transform.SizeComponent) {
	if maxSize, ok := t.maxSizeArray.Get(entity); ok {
		if maxSize.Size[0] != 0 && size.Size[0] > maxSize.Size[0] {
			size.Size[0] = maxSize.Size[0]
		}
		if maxSize.Size[1] != 0 && size.Size[1] > maxSize.Size[1] {
			size.Size[1] = maxSize.Size[1]
		}
		if maxSize.Size[2] != 0 && size.Size[2] > maxSize.Size[2] {
			size.Size[2] = maxSize.Size[2]
		}
	}
	if minSize, ok := t.minSizeArray.Get(entity); ok {
		if minSize.Size[0] != 0 && size.Size[0] < minSize.Size[0] {
			size.Size[0] = minSize.Size[0]
		}
		if minSize.Size[1] != 0 && size.Size[1] < minSize.Size[1] {
			size.Size[1] = minSize.Size[1]
		}
		if minSize.Size[2] != 0 && size.Size[2] < minSize.Size[2] {
			size.Size[2] = minSize.Size[2]
		}
	}
}

// its used by full ApplyAspectRatio method.
// ignores min and max size.
// also concludes that aspect ratio is verified.
func (t tool) getAspectRatio(size transform.SizeComponent, ratio transform.AspectRatioComponent) transform.SizeComponent {
	var base float32
	switch ratio.PrimaryAxis {
	case transform.PrimaryAxisX:
		base = size.Size[0] / ratio.AspectRatio[0]
	case transform.PrimaryAxisY:
		base = size.Size[1] / ratio.AspectRatio[1]
	case transform.PrimaryAxisZ:
		base = size.Size[2] / ratio.AspectRatio[2]
	default:
		return size
	}
	if ratio.AspectRatio[0] != 0 {
		size.Size[0] = base * ratio.AspectRatio[0]
	}
	if ratio.AspectRatio[1] != 0 {
		size.Size[1] = base * ratio.AspectRatio[1]
	}
	if ratio.AspectRatio[2] != 0 {
		size.Size[2] = base * ratio.AspectRatio[2]
	}
	return size
}

// integrates min and max size
func (t tool) ApplyAspectRatio(entity ecs.EntityID, size *transform.SizeComponent) {
	ratio, ok := t.aspectRatioArray.Get(entity)
	if !ok {
		return
	}
	if ratio.PrimaryAxis == 0 || ratio.PrimaryAxis > 3 || ratio.AspectRatio[ratio.PrimaryAxis-1] == 0 {
		return
	}
	sizeRatio := t.getAspectRatio(*size, ratio)
	if maxSize, ok := t.maxSizeArray.Get(entity); ok {
		for i := 0; i < 3; i++ {
			if maxSize.Size[i] == 0 {
				maxSize.Size[i] = sizeRatio.Size[i]
			}
		}
		maxSizeRatio := ratio

		primaryAxisIndex := -1
		for i := 0; i < 3; i++ {
			if maxSize.Size[i] != 0 && maxSizeRatio.AspectRatio[i] != 0 &&
				(primaryAxisIndex == -1 || maxSizeRatio.AspectRatio[primaryAxisIndex] > maxSizeRatio.AspectRatio[i]) {
				primaryAxisIndex = i
			}
		}
		if primaryAxisIndex != -1 {
			maxSizeRatio.PrimaryAxis = transform.PrimaryAxis(primaryAxisIndex + 1) // +1 because 0 isn't an axis
			maxSize = transform.MaxSizeComponent(
				t.getAspectRatio(transform.SizeComponent(maxSize), maxSizeRatio),
			)
			for i := 0; i < 3; i++ {
				if maxSize.Size[i] >= sizeRatio.Size[i] {
					continue
				}
				sizeRatio = transform.SizeComponent(maxSize)
				break
			}
		}
	}
	if minSize, ok := t.minSizeArray.Get(entity); ok {
		minSizeRatio := ratio

		primaryAxisIndex := -1
		for i := 0; i < 3; i++ {
			if minSize.Size[i] != 0 && minSizeRatio.AspectRatio[i] != 0 &&
				(primaryAxisIndex == -1 || minSizeRatio.AspectRatio[primaryAxisIndex] < minSizeRatio.AspectRatio[i]) {
				primaryAxisIndex = i
			}
		}
		if primaryAxisIndex != -1 {
			minSizeRatio.PrimaryAxis = transform.PrimaryAxis(primaryAxisIndex + 1) // +1 because 0 isn't an axis
			minSize = transform.MinSizeComponent(
				t.getAspectRatio(transform.SizeComponent(minSize), minSizeRatio),
			)
			for i := 0; i < 3; i++ {
				if minSize.Size[i] <= sizeRatio.Size[i] {
					continue
				}
				sizeRatio = transform.SizeComponent(minSize)
				break
			}
		}
	}
	*size = sizeRatio
}

//

func (t tool) CalculateAbsolutePos(entity ecs.EntityID) transform.AbsolutePosComponent {
	pos, _ := t.posArray.Get(entity)
	relativeToParentPos := t.GetRelativeParentPos(entity)

	pos.Pos = pos.Pos.
		Add(relativeToParentPos).
		Add(t.GetPivotPos(entity))

	return transform.AbsolutePosComponent(pos)
}
func (t tool) CalculateAbsoluteRot(entity ecs.EntityID) transform.AbsoluteRotationComponent {
	rot, ok := t.rotationArray.Get(entity)
	if !ok {
		rot = t.defaultRot
	}

	rot.Rotation = rot.Rotation.
		Mul(t.GetRelativeParentRotation(entity))

	return transform.AbsoluteRotationComponent(rot)
}
func (t tool) CalculateAbsoluteSize(entity ecs.EntityID) transform.AbsoluteSizeComponent {
	size, ok := t.sizeArray.Get(entity)
	if !ok {
		size = t.defaultSize
	}

	relativeParentSize := t.GetRelativeParentSize(entity)
	size.Size = mgl32.Vec3{
		size.Size[0] * relativeParentSize[0],
		size.Size[1] * relativeParentSize[1],
		size.Size[2] * relativeParentSize[2],
	}
	t.ApplyAspectRatio(entity, &size)
	t.ApplyMinMaxSize(entity, &size)

	return transform.AbsoluteSizeComponent(size)
}
