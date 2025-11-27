package transformtool

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"
)

type tool struct {
	logger logger.Logger

	world                ecs.World
	hierarchyTransaction hierarchy.Transaction

	defaultPos         transform.PosComponent
	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent

	parentArray           ecs.ComponentsArray[hierarchy.ParentComponent]
	posArray              ecs.ComponentsArray[transform.PosComponent]
	rotationArray         ecs.ComponentsArray[transform.RotationComponent]
	sizeArray             ecs.ComponentsArray[transform.SizeComponent]
	pivotPointArray       ecs.ComponentsArray[transform.PivotPointComponent]
	parentMaskArray       ecs.ComponentsArray[transform.ParentComponent]
	parentPivotPointArray ecs.ComponentsArray[transform.ParentPivotPointComponent]
}

func NewTransformTool(
	logger logger.Logger,
	hierarchyToolFactory ecs.ToolFactory[hierarchy.Tool],
	defaultPos transform.PosComponent,
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) ecs.ToolFactory[transform.Tool] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) transform.Tool {
		mutex.Lock()
		defer mutex.Unlock()

		if tool, err := ecs.GetGlobal[tool](w); err == nil {
			return tool
		}
		tool := tool{
			logger,
			w,
			hierarchyToolFactory.Build(w).Transaction(),
			defaultPos,
			defaultRot,
			defaultSize,
			defaultPivot,
			defaultParentPivot,
			ecs.GetComponentsArray[hierarchy.ParentComponent](w),
			ecs.GetComponentsArray[transform.PosComponent](w),
			ecs.GetComponentsArray[transform.RotationComponent](w),
			ecs.GetComponentsArray[transform.SizeComponent](w),
			ecs.GetComponentsArray[transform.PivotPointComponent](w),
			ecs.GetComponentsArray[transform.ParentComponent](w),
			ecs.GetComponentsArray[transform.ParentPivotPointComponent](w),
		}
		w.SaveGlobal(tool)
		tool.Init()
		return tool
	})
}

func (tool tool) Transaction() transform.Transaction {
	return newTransformTransaction(tool)
}

func (tool tool) Query(b ecs.LiveQueryBuilder) ecs.LiveQueryBuilder {
	return b.Track(
		ecs.GetComponentType(transform.PosComponent{}),
		ecs.GetComponentType(transform.RotationComponent{}),
		ecs.GetComponentType(transform.SizeComponent{}),
		ecs.GetComponentType(transform.PivotPointComponent{}),
		ecs.GetComponentType(transform.ParentComponent{}),
		ecs.GetComponentType(transform.ParentPivotPointComponent{}),
	)
}
