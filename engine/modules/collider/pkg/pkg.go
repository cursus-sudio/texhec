package colliderpkg

import (
	"engine/modules/collider"
	"engine/modules/collider/internal/collisions"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg { return pkg{} }

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) collider.Service {
		return collisions.NewService(c, 1000)
	})
}
