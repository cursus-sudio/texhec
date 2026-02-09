package transformservice

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type service struct {
	Logger logger.Logger `inject:"1"`

	World     ecs.World         `inject:"1"`
	Hierarchy hierarchy.Service `inject:"1"`

	DirtySet ecs.DirtySet

	AbsolutePosArray      ecs.ComponentsArray[transform.AbsolutePosComponent]
	AbsoluteSizeArray     ecs.ComponentsArray[transform.AbsoluteSizeComponent]
	AbsoluteRotationArray ecs.ComponentsArray[transform.AbsoluteRotationComponent]

	AbsolutePosWrapper      ecs.ComponentsArray[transform.AbsolutePosComponent]
	AbsoluteSizeWrapper     ecs.ComponentsArray[transform.AbsoluteSizeComponent]
	AbsoluteRotationWrapper ecs.ComponentsArray[transform.AbsoluteRotationComponent]

	PosArray      ecs.ComponentsArray[transform.PosComponent]
	RotationArray ecs.ComponentsArray[transform.RotationComponent]
	SizeArray     ecs.ComponentsArray[transform.SizeComponent]

	MaxSizeArray     ecs.ComponentsArray[transform.MaxSizeComponent]
	MinSizeArray     ecs.ComponentsArray[transform.MinSizeComponent]
	AspectRatioArray ecs.ComponentsArray[transform.AspectRatioComponent]

	PivotPointArray       ecs.ComponentsArray[transform.PivotPointComponent]
	ParentMaskArray       ecs.ComponentsArray[transform.ParentComponent]
	ParentPivotPointArray ecs.ComponentsArray[transform.ParentPivotPointComponent]

	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent
}

func NewService(
	c ioc.Dic,
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) transform.Service {
	s := ioc.GetServices[*service](c)

	s.DirtySet = ecs.NewDirtySet()

	s.AbsolutePosArray = ecs.GetComponentsArray[transform.AbsolutePosComponent](s.World)
	s.AbsoluteSizeArray = ecs.GetComponentsArray[transform.AbsoluteSizeComponent](s.World)
	s.AbsoluteRotationArray = ecs.GetComponentsArray[transform.AbsoluteRotationComponent](s.World)

	s.AbsolutePosWrapper = &absolutePosArray{s, s.AbsolutePosArray}
	s.AbsoluteSizeWrapper = &absoluteSizeArray{s, s.AbsoluteSizeArray}
	s.AbsoluteRotationWrapper = &absoluteRotationArray{s, s.AbsoluteRotationArray}

	s.PosArray = ecs.GetComponentsArray[transform.PosComponent](s.World)
	s.SizeArray = ecs.GetComponentsArray[transform.SizeComponent](s.World)
	s.RotationArray = ecs.GetComponentsArray[transform.RotationComponent](s.World)

	s.MaxSizeArray = ecs.GetComponentsArray[transform.MaxSizeComponent](s.World)
	s.MinSizeArray = ecs.GetComponentsArray[transform.MinSizeComponent](s.World)
	s.AspectRatioArray = ecs.GetComponentsArray[transform.AspectRatioComponent](s.World)

	s.PivotPointArray = ecs.GetComponentsArray[transform.PivotPointComponent](s.World)
	s.ParentMaskArray = ecs.GetComponentsArray[transform.ParentComponent](s.World)
	s.ParentPivotPointArray = ecs.GetComponentsArray[transform.ParentPivotPointComponent](s.World)

	s.defaultRot = defaultRot
	s.defaultSize = defaultSize
	s.defaultPivot = defaultPivot
	s.defaultParentPivot = defaultParentPivot

	s.Init()
	return s

}

func (t *service) BeforeGet() {
	entities := t.DirtySet.Get()
	if len(entities) == 0 {
		return
	}
	children := []ecs.EntityID{}

	saves := []save{}

	for len(entities) != 0 || len(children) != 0 {
		if len(entities) == 0 {
			for _, save := range saves {
				t.AbsolutePosArray.Set(save.entity, save.pos)
				t.AbsoluteRotationArray.Set(save.entity, save.rot)
				t.AbsoluteSizeArray.Set(save.entity, save.size)
			}
			t.DirtySet.Clear()

			entities = children
			children = nil
			saves = nil
		}
		entity := entities[0]
		entities = entities[1:]
		if !t.World.EntityExists(entity) {
			continue
		}

		pos, rot, size := t.CalculateAbsolute(entity)
		save := save{
			entity: entity,
			pos:    pos, rot: rot, size: size,
		}

		saves = append(saves, save)

		for _, child := range t.Hierarchy.Children(entity).GetIndices() {
			comparedMask := transform.RelativePos | transform.RelativeRotation | transform.RelativeSizeXYZ
			mask, _ := t.ParentMaskArray.Get(child)
			if mask.RelativeMask&comparedMask == 0 {
				continue
			}
			children = append(children, child)
		}
	}

	for _, save := range saves {
		t.AbsolutePosArray.Set(save.entity, save.pos)
		t.AbsoluteRotationArray.Set(save.entity, save.rot)
		t.AbsoluteSizeArray.Set(save.entity, save.size)
	}
	t.DirtySet.Clear()
}

func (t *service) AbsolutePos() ecs.ComponentsArray[transform.AbsolutePosComponent] {
	return t.AbsolutePosWrapper
}
func (t *service) AbsoluteRotation() ecs.ComponentsArray[transform.AbsoluteRotationComponent] {
	return t.AbsoluteRotationWrapper
}
func (t *service) AbsoluteSize() ecs.ComponentsArray[transform.AbsoluteSizeComponent] {
	return t.AbsoluteSizeWrapper
}
func (t *service) Pos() ecs.ComponentsArray[transform.PosComponent] {
	return t.PosArray
}
func (t *service) Rotation() ecs.ComponentsArray[transform.RotationComponent] {
	return t.RotationArray
}
func (t *service) Size() ecs.ComponentsArray[transform.SizeComponent] {
	return t.SizeArray
}
func (t *service) MaxSize() ecs.ComponentsArray[transform.MaxSizeComponent] {
	return t.MaxSizeArray
}
func (t *service) MinSize() ecs.ComponentsArray[transform.MinSizeComponent] {
	return t.MinSizeArray
}
func (t *service) AspectRatio() ecs.ComponentsArray[transform.AspectRatioComponent] {
	return t.AspectRatioArray
}
func (t *service) PivotPoint() ecs.ComponentsArray[transform.PivotPointComponent] {
	return t.PivotPointArray
}
func (t *service) Parent() ecs.ComponentsArray[transform.ParentComponent] {
	return t.ParentMaskArray
}
func (t *service) ParentPivotPoint() ecs.ComponentsArray[transform.ParentPivotPointComponent] {
	return t.ParentPivotPointArray
}

func (t *service) Mat4(entity ecs.EntityID) mgl32.Mat4 {
	pos, ok := t.AbsolutePosArray.Get(entity)
	if !ok {
		pos.Pos = mgl32.Vec3{0, 0, 0}
	}
	rot, ok := t.AbsoluteRotationArray.Get(entity)
	if !ok {
		rot.Rotation = mgl32.QuatIdent()
	}
	size, ok := t.AbsoluteSizeArray.Get(entity)
	if !ok {
		size.Size = mgl32.Vec3{1, 1, 1}
	}

	translation := mgl32.Translate3D(pos.Pos.X(), pos.Pos.Y(), pos.Pos.Z())
	rotation := rot.Rotation.Mat4()
	scale := mgl32.Scale3D(size.Size.X()/2, size.Size.Y()/2, size.Size.Z()/2)
	return translation.Mul4(rotation).Mul4(scale)
}

func (t *service) AddDirtySet(set ecs.DirtySet) {
	t.AbsolutePosArray.AddDirtySet(set)
	t.AbsoluteRotationArray.AddDirtySet(set)
	t.AbsoluteSizeArray.AddDirtySet(set)
}
