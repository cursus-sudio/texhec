package transformpkg

import (
	"frontend/modules/transform"
	"frontend/modules/transform/internal/transformtool"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	defaultPos         transform.PosComponent
	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent
}

func Package(
	defaultPos transform.PosComponent,
	defaultRot transform.RotationComponent,
	defaultSize transform.SizeComponent,
	defaultPivot transform.PivotPointComponent,
	defaultParentPivot transform.ParentPivotPointComponent,
) ioc.Pkg {
	return pkg{
		defaultPos,
		defaultRot,
		defaultSize,
		defaultPivot,
		defaultParentPivot,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[transform.TransformTool] {
		return transformtool.NewTransformTool(
			ioc.Get[logger.Logger](c),
			pkg.defaultPos,
			pkg.defaultRot,
			pkg.defaultSize,
			pkg.defaultPivot,
			pkg.defaultParentPivot,
		)
	})
}
