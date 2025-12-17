package transformtool

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type tool struct {
	logger logger.Logger

	world     ecs.World
	dirtySet  ecs.DirtySet
	hierarchy hierarchy.Interface

	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent

	absolutePosArray      ecs.ComponentsArray[transform.AbsolutePosComponent]
	absoluteRotationArray ecs.ComponentsArray[transform.AbsoluteRotationComponent]
	absoluteSizeArray     ecs.ComponentsArray[transform.AbsoluteSizeComponent]

	hierarchyArray        ecs.ComponentsArray[hierarchy.Component]
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

func NewTransformTool(
	logger logger.Logger,
	hierarchyToolFactory ecs.ToolFactory[hierarchy.HierarchyTool],
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) ecs.ToolFactory[transform.TransformTool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) transform.TransformTool {
		mutex.Lock()
		defer mutex.Unlock()

		if tool, ok := ecs.GetGlobal[tool](w); ok {
			return tool
		}
		tool := tool{
			logger,
			w,
			ecs.NewDirtySet(),
			hierarchyToolFactory.Build(w).Hierarchy(),
			defaultRot,
			defaultSize,
			defaultPivot,
			defaultParentPivot,
			ecs.GetComponentsArray[transform.AbsolutePosComponent](w),
			ecs.GetComponentsArray[transform.AbsoluteRotationComponent](w),
			ecs.GetComponentsArray[transform.AbsoluteSizeComponent](w),
			ecs.GetComponentsArray[hierarchy.Component](w),
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
		w.SaveGlobal(tool)
		tool.Init()
		return tool
	})
}

func (t tool) BeforeGet() {
	entities := t.dirtySet.Get()
	if len(entities) == 0 {
		return
	}
	children := []ecs.EntityID{}

	saves := []save{}

	for len(entities) != 0 || len(children) != 0 {
		if len(entities) == 0 {
			entities = children
			for _, save := range saves {
				t.absolutePosArray.Set(save.entity, save.pos)
				t.absoluteRotationArray.Set(save.entity, save.rot)
				t.absoluteSizeArray.Set(save.entity, save.size)
			}
			t.dirtySet.Clear()

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
			mask, ok := t.parentMaskArray.Get(child)
			if !ok || mask.RelativeMask&comparedMask == 0 {
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

func (t tool) SetAbsolutePos(entity ecs.EntityID, absolutePos transform.AbsolutePosComponent) {
	pos := transform.NewPos(absolutePos.Pos.
		Sub(t.GetRelativeParentPos(entity)).
		Sub(t.GetPivotPos(entity)).Elem())

	t.posArray.Set(entity, pos)
}
func (t tool) SetAbsoluteRotation(entity ecs.EntityID, absoluteRot transform.AbsoluteRotationComponent) {
	rot := transform.NewRotation(absoluteRot.Rotation.
		Mul(t.GetRelativeParentRotation(entity).Inverse()))

	t.rotationArray.Set(entity, rot)
}
func (t tool) SetAbsoluteSize(entity ecs.EntityID, absoluteSize transform.AbsoluteSizeComponent) {
	parentSize := t.GetRelativeParentSize(entity)
	size := transform.NewSize(
		absoluteSize.Size[0]/parentSize[0],
		absoluteSize.Size[1]/parentSize[1],
		absoluteSize.Size[2]/parentSize[2],
	)

	t.sizeArray.Set(entity, size)
}

func (t tool) Transform() transform.Interface { return t }

func (t tool) AbsolutePos() ecs.ComponentsArray[transform.AbsolutePosComponent] {
	return t.absolutePosArray
}
func (t tool) AbsoluteRotation() ecs.ComponentsArray[transform.AbsoluteRotationComponent] {
	return t.absoluteRotationArray
}
func (t tool) AbsoluteSize() ecs.ComponentsArray[transform.AbsoluteSizeComponent] {
	return t.absoluteSizeArray
}
func (t tool) Pos() ecs.ComponentsArray[transform.PosComponent] {
	return t.posArray
}
func (t tool) Rotation() ecs.ComponentsArray[transform.RotationComponent] {
	return t.rotationArray
}
func (t tool) Size() ecs.ComponentsArray[transform.SizeComponent] {
	return t.sizeArray
}
func (t tool) MaxSize() ecs.ComponentsArray[transform.MaxSizeComponent] {
	return t.maxSizeArray
}
func (t tool) MinSize() ecs.ComponentsArray[transform.MinSizeComponent] {
	return t.minSizeArray
}
func (t tool) AspectRatio() ecs.ComponentsArray[transform.AspectRatioComponent] {
	return t.aspectRatioArray
}
func (t tool) PivotPoint() ecs.ComponentsArray[transform.PivotPointComponent] {
	return t.pivotPointArray
}
func (t tool) Parent() ecs.ComponentsArray[transform.ParentComponent] {
	return t.parentMaskArray
}
func (t tool) ParentPivotPoint() ecs.ComponentsArray[transform.ParentPivotPointComponent] {
	return t.parentPivotPointArray
}

func (t tool) Mat4(entity ecs.EntityID) mgl32.Mat4 {
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

func (t tool) AddDirtySet(set ecs.DirtySet) {
	t.absolutePosArray.AddDirtySet(set)
	t.absoluteRotationArray.AddDirtySet(set)
	t.absoluteSizeArray.AddDirtySet(set)
}
