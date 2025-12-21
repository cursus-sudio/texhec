package recordpkg

import (
	"engine/modules/record"
	"engine/modules/record/internal/recordimpl"
	"engine/modules/uuid"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

// this is parent configuration.
// it should have all used recorded components in any configuration.
// note: each new recorded component in configuration adds new BeforeGet to this type
// so do not add it automatically to everyhing because it can result in performance loss
func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) record.ToolFactory {
		return recordimpl.NewToolFactory(
			ioc.Get[uuid.ToolFactory](c),
			ioc.Get[logger.Logger](c),
		)
	})
}
