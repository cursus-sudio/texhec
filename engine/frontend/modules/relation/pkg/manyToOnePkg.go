package relationpkg

import (
	"frontend/modules/relation"
	"frontend/modules/relation/internal/manytoone"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type manyToOnePkg[Component any] struct {
	componentParent func(Component) ecs.EntityID
}

func ManyToOnePackage[Component any](componentParent func(Component) ecs.EntityID) ioc.Pkg {
	return manyToOnePkg[Component]{componentParent}
}

func (pkg manyToOnePkg[Component]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[relation.EntityToEntitiesTool[Component]] {
		return manytoone.NewManyToOneFactory(pkg.componentParent)
	})
}
