package uuidpkg

import (
	"engine/modules/relation"
	relationpkg "engine/modules/relation/pkg"
	uuid "engine/modules/uuid"
	"engine/modules/uuid/internal"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) uuid.Factory { return internal.NewFactory() })
	relationpkg.MapRelationPackage(
		func(w ecs.World) ecs.LiveQuery {
			return w.Query().
				Require(uuid.Component{}).
				Build()
		},
		func(w ecs.World) func(entity ecs.EntityID) (indexType uuid.UUID, ok bool) {
			uniqueArray := ecs.GetComponentsArray[uuid.Component](w)
			return func(entity ecs.EntityID) (indexType uuid.UUID, ok bool) {
				component, err := uniqueArray.GetComponent(entity)
				if err != nil {
					return uuid.UUID{}, false
				}
				return component.ID, true
			}
		},
	).Register(b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[uuid.Tool] {
		return internal.NewToolFactory(
			ioc.Get[ecs.ToolFactory[relation.EntityToKeyTool[uuid.UUID]]](c),
			ioc.Get[uuid.Factory](c),
		)
	})
}
