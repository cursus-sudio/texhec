package transformtool

import (
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type tool struct {
	logger logger.Logger

	world    transform.World
	dirtySet ecs.DirtySet

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

func NewTransformTool(
	logger logger.Logger,
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) transform.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w transform.World) transform.TransformTool {
		mutex.Lock()
		defer mutex.Unlock()

		if tool, ok := ecs.GetGlobal[tool](w); ok {
			return tool
		}
		tool := &tool{
			logger,
			w,
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
		// tool.Interface = tool

		tool.absolutePosWrapper = &absolutePosArray{tool, tool.absolutePosArray}
		tool.absoluteSizeWrapper = &absoluteSizeArray{tool, tool.absoluteSizeArray}
		tool.absoluteRotationWrapper = &absoluteRotationArray{tool, tool.absoluteRotationArray}

		w.SaveGlobal(tool)
		tool.Init()
		return tool
	})
}

func (t *tool) BeforeGet() {
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

		for _, child := range t.world.Hierarchy().Children(entity).GetIndices() {
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

func (t *tool) Transform() transform.Interface { return t }

func (t *tool) AbsolutePos() ecs.ComponentsArray[transform.AbsolutePosComponent] {
	return t.absolutePosWrapper
}
func (t *tool) AbsoluteRotation() ecs.ComponentsArray[transform.AbsoluteRotationComponent] {
	return t.absoluteRotationWrapper
}
func (t *tool) AbsoluteSize() ecs.ComponentsArray[transform.AbsoluteSizeComponent] {
	return t.absoluteSizeWrapper
}
func (t *tool) Pos() ecs.ComponentsArray[transform.PosComponent] {
	return t.posArray
}
func (t *tool) Rotation() ecs.ComponentsArray[transform.RotationComponent] {
	return t.rotationArray
}
func (t *tool) Size() ecs.ComponentsArray[transform.SizeComponent] {
	return t.sizeArray
}
func (t *tool) MaxSize() ecs.ComponentsArray[transform.MaxSizeComponent] {
	return t.maxSizeArray
}
func (t *tool) MinSize() ecs.ComponentsArray[transform.MinSizeComponent] {
	return t.minSizeArray
}
func (t *tool) AspectRatio() ecs.ComponentsArray[transform.AspectRatioComponent] {
	return t.aspectRatioArray
}
func (t *tool) PivotPoint() ecs.ComponentsArray[transform.PivotPointComponent] {
	return t.pivotPointArray
}
func (t *tool) Parent() ecs.ComponentsArray[transform.ParentComponent] {
	return t.parentMaskArray
}
func (t *tool) ParentPivotPoint() ecs.ComponentsArray[transform.ParentPivotPointComponent] {
	return t.parentPivotPointArray
}

func (t *tool) Mat4(entity ecs.EntityID) mgl32.Mat4 {
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

func (t *tool) AddDirtySet(set ecs.DirtySet) {
	t.absolutePosArray.AddDirtySet(set)
	t.absoluteRotationArray.AddDirtySet(set)
	t.absoluteSizeArray.AddDirtySet(set)
}
