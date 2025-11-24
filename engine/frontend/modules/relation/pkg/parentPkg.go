package relationpkg

import (
	"frontend/modules/relation"
	"frontend/modules/relation/internal/parent"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type parentPkg[Component any] struct {
	componentParent func(Component) ecs.EntityID
}

func ParentPackage[Component any](componentParent func(Component) ecs.EntityID) ioc.Pkg {
	return parentPkg[Component]{componentParent}
}

func (pkg parentPkg[Component]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[relation.ParentTool[Component]] {
		return parent.NewParentToolFactory(pkg.componentParent)
	})
}
