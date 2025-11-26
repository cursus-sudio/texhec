package transformtool

import (
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

func (t object) GetRelativeParentPos() mgl32.Vec3 {
	parent, err := t.parent.Get()
	if err != nil {
		return mgl32.Vec3{}
	}
	parentMask, err := t.parentMask.Get()
	if err != nil || parentMask.RelativeMask&transform.RelativePos == 0 {
		return mgl32.Vec3{}
	}
	parentTransform := t.GetObject(parent.Parent)
	parentPos, err := parentTransform.AbsolutePos().Get()
	if err != nil {
		parentPos = t.defaultPos
	}
	parentSize, err := parentTransform.AbsoluteSize().Get()
	if err != nil {
		parentSize = t.defaultSize
	}
	parentPivot, err := t.parentPivotPoint.Get()
	if err != nil {
		parentPivot = t.defaultParentPivot
	}
	parentPivot.Point = parentPivot.Point.Sub(t.defaultParentPivot.Point)
	return parentPos.Pos.Add(mgl32.Vec3{
		parentSize.Size[0] * parentPivot.Point[0],
		parentSize.Size[1] * parentPivot.Point[1],
		parentSize.Size[2] * parentPivot.Point[2],
	})
}

func (t object) GetPivotPos() mgl32.Vec3 {
	pivot, err := t.pivotPoint.Get()
	if err != nil {
		return mgl32.Vec3{}
	}
	size, err := t.absoluteSize.Get()
	if err != nil {
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

func (t object) GetRelativeParentRotation() mgl32.Quat {
	parent, err := t.parent.Get()
	if err != nil {
		return mgl32.QuatIdent()
	}
	parentMask, err := t.parentMask.Get()
	if err != nil || parentMask.RelativeMask&transform.RelativeRotation == 0 {
		return mgl32.QuatIdent()
	}
	parentTransform := t.GetObject(parent.Parent)
	parentRot, err := parentTransform.AbsoluteRotation().Get()
	if err != nil {
		return mgl32.QuatIdent()
	}
	return parentRot.Rotation
}

//

func (t object) GetRelativeParentSize() mgl32.Vec3 {
	parent, err := t.parent.Get()
	if err != nil {
		return mgl32.Vec3{1, 1, 1}
	}
	parentMask, err := t.parentMask.Get()
	if err != nil || parentMask.RelativeMask&transform.RelativeSize == 0 {
		return mgl32.Vec3{1, 1, 1}
	}
	parentTransform := t.GetObject(parent.Parent)
	parentSize, err := parentTransform.AbsoluteSize().Get()
	if err != nil {
		return mgl32.Vec3{1, 1, 1}
	}
	return parentSize.Size
}

//

func (t *object) Init() {
	t.absolutePos = ecs.NewEntityComponent(
		func() (transform.PosComponent, error) {
			pos, err := t.pos.Get()
			if err != nil {
				pos = t.defaultPos
			}
			relativeToParentPos := t.GetRelativeParentPos()

			pos.Pos = pos.Pos.
				Add(relativeToParentPos).
				Add(t.GetPivotPos())

			return pos, nil
		},
		func(absolutePos transform.PosComponent) {
			pos, err := t.pos.Get()
			if err != nil {
				pos.Pos = t.defaultPos.Pos
			}
			relativeToParentPos := t.GetRelativeParentPos()

			pos.Pos = absolutePos.Pos.
				Sub(relativeToParentPos).
				Sub(t.GetPivotPos())

			t.pos.Set(pos)
		},
		t.pos.Remove,
	)
	t.absoluteRot = ecs.NewEntityComponent(
		func() (transform.RotationComponent, error) {
			rot, err := t.rot.Get()
			if err != nil {
				rot = t.defaultRot
			}

			rot.Rotation = rot.Rotation.
				Mul(t.GetRelativeParentRotation())

			return rot, nil
		},
		func(absoluteRot transform.RotationComponent) {
			rot, err := t.rot.Get()
			if err != nil {
				rot = t.defaultRot
			}

			rot.Rotation = absoluteRot.Rotation.
				Mul(t.GetRelativeParentRotation().Inverse())

			t.rot.Set(rot)
		},
		t.rot.Remove,
	)
	t.absoluteSize = ecs.NewEntityComponent(
		func() (transform.SizeComponent, error) {
			size, err := t.size.Get()
			if err != nil {
				size = t.defaultSize
			}

			relativeParentSize := t.GetRelativeParentSize()
			size.Size = mgl32.Vec3{
				size.Size[0] * relativeParentSize[0],
				size.Size[1] * relativeParentSize[1],
				size.Size[2] * relativeParentSize[2],
			}

			return size, nil
		},
		func(absoluteSize transform.SizeComponent) {
			size, err := t.size.Get()
			if err != nil {
				size = t.defaultSize
			}

			size.Size = absoluteSize.Size.
				Sub(t.GetRelativeParentSize())

			t.size.Set(size)
		},
		t.size.Remove,
	)
}
