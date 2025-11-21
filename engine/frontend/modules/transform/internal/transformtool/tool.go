package transformtool

import (
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"
)

type transformTool struct {
	logger logger.Logger

	world ecs.World

	defaultPos         transform.PosComponent
	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent

	posArray              ecs.ComponentsArray[transform.PosComponent]
	rotationArray         ecs.ComponentsArray[transform.RotationComponent]
	sizeArray             ecs.ComponentsArray[transform.SizeComponent]
	pivotPointArray       ecs.ComponentsArray[transform.PivotPointComponent]
	parentArray           ecs.ComponentsArray[transform.ParentComponent]
	parentPivotPointArray ecs.ComponentsArray[transform.ParentPivotPointComponent]
}

func NewTransformTool(
	logger logger.Logger,
	defaultPos transform.PosComponent,
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) ecs.ToolFactory[transform.TransformTool] {
	return ecs.NewToolFactory(func(w ecs.World) transform.TransformTool {
		return transformTool{
			logger,
			w,
			defaultPos,
			defaultRot,
			defaultSize,
			defaultPivot,
			defaultParentPivot,
			ecs.GetComponentsArray[transform.PosComponent](w.Components()),
			ecs.GetComponentsArray[transform.RotationComponent](w.Components()),
			ecs.GetComponentsArray[transform.SizeComponent](w.Components()),
			ecs.GetComponentsArray[transform.PivotPointComponent](w.Components()),
			ecs.GetComponentsArray[transform.ParentComponent](w.Components()),
			ecs.GetComponentsArray[transform.ParentPivotPointComponent](w.Components()),
		}
	})
}

func (tool transformTool) Transaction() transform.TransformTransaction {
	return newTransformTransaction(tool)
}

func (tool transformTool) Query(b ecs.LiveQueryBuilder) ecs.LiveQueryBuilder {
	return b.Track(
		ecs.GetComponentType(transform.PosComponent{}),
		ecs.GetComponentType(transform.RotationComponent{}),
		ecs.GetComponentType(transform.SizeComponent{}),
		ecs.GetComponentType(transform.PivotPointComponent{}),
		ecs.GetComponentType(transform.ParentComponent{}),
		ecs.GetComponentType(transform.ParentPivotPointComponent{}),
	)
}
