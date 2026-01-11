package definitionpkg

import (
	"core/modules/definition"
	"core/modules/definition/internal"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			Register(definition.DefinitionComponent{}).
			Register(definition.DefinitionLinkComponent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) definition.ToolFactory {
		return internal.NewToolFactory()
	})
}
