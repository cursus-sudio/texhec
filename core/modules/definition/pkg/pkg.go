package definitionpkg

import (
	"core/modules/definition"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			Register(definition.DefinitionComponent{}).
			Register(definition.DefinitionLinkComponent{})
	})
}
