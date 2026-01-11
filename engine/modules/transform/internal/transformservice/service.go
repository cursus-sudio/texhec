package transformservice

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
)

type service struct {
	logger logger.Logger

	world     ecs.World
	hierarchy hierarchy.Service
	dirtySet  ecs.DirtySet

	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent

	absolutePosArray      ecs.ComponentsArray[transform.AbsolutePosComponent]
	absoluteSizeArray     ecs.ComponentsArray[transform.AbsoluteSizeComponent]
	absoluteRotationArray ecs.ComponentsArray[transform.AbsoluteRotationComponent]

	absolutePosWrapper      ecs.ComponentsArray[transform.AbsolutePosComponent]
	absoluteSizeWrapper     ecs.ComponentsArray[transform.AbsoluteSizeComponent]
	absoluteRotationWrapper ecs.ComponentsArray[transform.AbsoluteRotationComponent]

	posArray              ecs.ComponentsArray[transform.PosComponent]
	rotationArray         ecs.ComponentsArray[transform.RotationComponent]
	sizeArray             ecs.ComponentsArray[transform.SizeComponent]
	maxSizeArray          ecs.ComponentsArray[transform.MaxSizeComponent]
	minSizeArray          ecs.ComponentsArray[transform.MinSizeComponent]
	aspectRatioArray      ecs.ComponentsArray[transform.AspectRatioComponent]
	pivotPointArray       ecs.ComponentsArray[transform.PivotPointComponent]
	parentMaskArray       ecs.ComponentsArray[transform.ParentComponent]
	parentPivotPointArray ecs.ComponentsArray[transform.ParentPivotPointComponent]
}

func NewService(
	w ecs.World,
	hierarchy hierarchy.Service,
	logger logger.Logger,
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) transform.Service {
	s := &service{
		logger,
		w,
		hierarchy,
		ecs.NewDirtySet(),
		defaultRot,
		defaultSize,
		defaultPivot,
		defaultParentPivot,
		ecs.GetComponentsArray[transform.AbsolutePosComponent](w),
		ecs.GetComponentsArray[transform.AbsoluteSizeComponent](w),
		ecs.GetComponentsArray[transform.AbsoluteRotationComponent](w),
		nil,
		nil,
		nil,
		ecs.GetComponentsArray[transform.PosComponent](w),
		ecs.GetComponentsArray[transform.RotationComponent](w),
		ecs.GetComponentsArray[transform.SizeComponent](w),
		ecs.GetComponentsArray[transform.MaxSizeComponent](w),
		ecs.GetComponentsArray[transform.MinSizeComponent](w),
		ecs.GetComponentsArray[transform.AspectRatioComponent](w),
		ecs.GetComponentsArray[transform.PivotPointComponent](w),
		ecs.GetComponentsArray[transform.ParentComponent](w),
		ecs.GetComponentsArray[transform.ParentPivotPointComponent](w),
	}

	s.absolutePosWrapper = &absolutePosArray{s, s.absolutePosArray}
	s.absoluteSizeWrapper = &absoluteSizeArray{s, s.absoluteSizeArray}
	s.absoluteRotationWrapper = &absoluteRotationArray{s, s.absoluteRotationArray}

	s.Init()
	return s

}

func (t *service) BeforeGet() {
	entities := t.dirtySet.Get()
	if len(entities) == 0 {
		return
	}
	children := []ecs.EntityID{}

	saves := []save{}

	for len(entities) != 0 || len(children) != 0 {
		if len(entities) == 0 {
			for _, save := range saves {
				t.absolutePosArray.Set(save.entity, save.pos)
				t.absoluteRotationArray.Set(save.entity, save.rot)
				t.absoluteSizeArray.Set(save.entity, save.size)
			}
			t.dirtySet.Clear()

			entities = children
			children = nil
			saves = nil
		}
		entity := entities[0]
		entities = entities[1:]

		pos, rot, size := t.CalculateAbsolute(entity)
		save := save{
			entity: entity,
			pos:    pos, rot: rot, size: size,
		}

		saves = append(saves, save)

		for _, child := range t.hierarchy.Children(entity).GetIndices() {
			comparedMask := transform.RelativePos | transform.RelativeRotation | transform.RelativeSizeXYZ
			mask, _ := t.parentMaskArray.Get(child)
			if mask.RelativeMask&comparedMask == 0 {
				continue
			}
			children = append(children, child)
		}
	}

	for _, save := range saves {
		t.absolutePosArray.Set(save.entity, save.pos)
		t.absoluteRotationArray.Set(save.entity, save.rot)
		t.absoluteSizeArray.Set(save.entity, save.size)
	}
	t.dirtySet.Clear()
}

func (t *service) AbsolutePos() ecs.ComponentsArray[transform.AbsolutePosComponent] {
	return t.absolutePosWrapper
}
func (t *service) AbsoluteRotation() ecs.ComponentsArray[transform.AbsoluteRotationComponent] {
	return t.absoluteRotationWrapper
}
func (t *service) AbsoluteSize() ecs.ComponentsArray[transform.AbsoluteSizeComponent] {
	return t.absoluteSizeWrapper
}
func (t *service) Pos() ecs.ComponentsArray[transform.PosComponent] {
	return t.posArray
}
func (t *service) Rotation() ecs.ComponentsArray[transform.RotationComponent] {
	return t.rotationArray
}
func (t *service) Size() ecs.ComponentsArray[transform.SizeComponent] {
	return t.sizeArray
}
func (t *service) MaxSize() ecs.ComponentsArray[transform.MaxSizeComponent] {
	return t.maxSizeArray
}
func (t *service) MinSize() ecs.ComponentsArray[transform.MinSizeComponent] {
	return t.minSizeArray
}
func (t *service) AspectRatio() ecs.ComponentsArray[transform.AspectRatioComponent] {
	return t.aspectRatioArray
}
func (t *service) PivotPoint() ecs.ComponentsArray[transform.PivotPointComponent] {
	return t.pivotPointArray
}
func (t *service) Parent() ecs.ComponentsArray[transform.ParentComponent] {
	return t.parentMaskArray
}
func (t *service) ParentPivotPoint() ecs.ComponentsArray[transform.ParentPivotPointComponent] {
	return t.parentPivotPointArray
}

func (t *service) Mat4(entity ecs.EntityID) mgl32.Mat4 {
	pos, ok := t.absolutePosArray.Get(entity)
	if !ok {
		pos.Pos = mgl32.Vec3{0, 0, 0}
	}
	rot, ok := t.absoluteRotationArray.Get(entity)
	if !ok {
		rot.Rotation = mgl32.QuatIdent()
	}
	size, ok := t.absoluteSizeArray.Get(entity)
	if !ok {
		size.Size = mgl32.Vec3{1, 1, 1}
	}

	translation := mgl32.Translate3D(pos.Pos.X(), pos.Pos.Y(), pos.Pos.Z())
	rotation := rot.Rotation.Mat4()
	scale := mgl32.Scale3D(size.Size.X()/2, size.Size.Y()/2, size.Size.Z()/2)
	return translation.Mul4(rotation).Mul4(scale)
}

func (t *service) AddDirtySet(set ecs.DirtySet) {
	t.absolutePosArray.AddDirtySet(set)
	t.absoluteRotationArray.AddDirtySet(set)
	t.absoluteSizeArray.AddDirtySet(set)
}
