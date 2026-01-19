package service

import (
	"engine/modules/hierarchy"
	"engine/modules/layout"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	Logger    logger.Logger     `inject:"1"`
	World     ecs.World         `inject:"1"`
	Hierarchy hierarchy.Service `inject:"1"`
	Transform transform.Service `inject:"1"`

	align         ecs.ComponentsArray[layout.AlignComponent]
	order         ecs.ComponentsArray[layout.OrderComponent]
	gap           ecs.ComponentsArray[layout.GapComponent]
	dirtyParents  ecs.DirtySet
	dirtyChildren ecs.DirtySet
}

func NewLayoutService(c ioc.Dic,
) layout.Service {
	t := ioc.GetServices[*service](c)
	t.align = ecs.GetComponentsArray[layout.AlignComponent](t.World)
	t.order = ecs.GetComponentsArray[layout.OrderComponent](t.World)
	t.gap = ecs.GetComponentsArray[layout.GapComponent](t.World)
	t.dirtyParents = ecs.NewDirtySet()
	t.dirtyChildren = ecs.NewDirtySet()
	t.Init()
	return t
}

func (t *service) Align() ecs.ComponentsArray[layout.AlignComponent] { return t.align }
func (t *service) Order() ecs.ComponentsArray[layout.OrderComponent] { return t.order }
func (t *service) Gap() ecs.ComponentsArray[layout.GapComponent]     { return t.gap }

//

func (t *service) Init() {
	// t.order.SetEmpty(layout.NewOrder(layout.OrderHorizontal))
	t.align.SetEmpty(layout.NewAlign(.5, .5))
	t.gap.SetEmpty(layout.NewGap(0))

	t.Transform.AbsolutePos().AddDependency(t.align)
	t.Transform.AbsolutePos().AddDependency(t.order)
	t.Transform.AbsolutePos().AddDependency(t.gap)

	t.align.AddDirtySet(t.dirtyParents)
	t.order.AddDirtySet(t.dirtyParents)
	t.gap.AddDirtySet(t.dirtyParents)
	t.Transform.AddDirtySet(t.dirtyParents)

	t.Transform.AddDirtySet(t.dirtyChildren)
	t.Hierarchy.Component().AddDirtySet(t.dirtyChildren)

	// before get trigger
	t.Transform.AbsolutePos().BeforeGet(t.BeforeGet)
	t.Transform.AbsoluteSize().BeforeGet(t.BeforeGet)
}

type save struct {
	entity      ecs.EntityID
	pos         transform.PosComponent
	pivot       transform.PivotPointComponent
	parentPivot transform.ParentPivotPointComponent
}

func (t *service) BeforeGet() {
	for _, child := range t.dirtyChildren.Get() {
		if parent, ok := t.Hierarchy.Parent(child); ok {
			t.dirtyParents.Dirty(parent)
		}
	}
	parents := t.dirtyParents.Get()
	if len(parents) == 0 {
		return
	}
	defer t.dirtyChildren.Clear()
	defer t.dirtyParents.Clear()

	saves := []save{}

	for _, parent := range parents {
		parentSaves := t.handleParentChildren(parent)
		saves = append(saves, parentSaves...)
	}

	for _, save := range saves {
		t.Transform.Pos().Set(save.entity, save.pos)
		t.Transform.PivotPoint().Set(save.entity, save.pivot)
		t.Transform.ParentPivotPoint().Set(save.entity, save.parentPivot)
	}
}

func (t *service) handleParentChildren(parent ecs.EntityID) []save {
	children := t.Hierarchy.Children(parent).GetIndices()
	if len(children) == 0 {
		return nil
	}
	order, ok := t.order.Get(parent)
	if !ok {
		return nil
	}
	saves := make([]save, 0, len(children))
	align, _ := t.align.Get(parent)
	gap, _ := t.gap.Get(parent)

	// including gaps
	var totalSize float32 = 0
	for _, child := range children {
		size, _ := t.Transform.AbsoluteSize().Get(child)
		totalSize += size.Size[order.Order] + gap.Gap
	}
	totalSize -= gap.Gap

	size, _ := t.Transform.AbsoluteSize().Get(parent)
	progress := totalSize - size.Size[order.Primary()]
	progress *= align.Primary

	for _, child := range children {
		// pos
		pos := transform.NewPos(0, 0, 1)
		pos.Pos[order.Primary()] = progress

		// pivot point
		pivot := transform.NewPivotPoint(.5, .5, .5)
		pivot.Point[order.Primary()] = 1
		pivot.Point[order.Secondary()] = align.Secondary

		// parent pivot
		parentPivot := transform.NewParentPivotPoint(.5, .5, .5)
		parentPivot.Point[order.Primary()] = 1
		parentPivot.Point[order.Secondary()] = align.Secondary

		save := save{
			child,
			pos,
			pivot,
			parentPivot,
		}
		saves = append(saves, save)

		// update progress
		size, _ := t.Transform.AbsoluteSize().Get(child)
		progress -= size.Size[order.Primary()] + gap.Gap
	}

	return saves
}
