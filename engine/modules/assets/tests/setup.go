package test

import (
	"engine/modules/assets"
	assetspkg "engine/modules/assets/pkg"
	"engine/modules/registry"
	registrypkg "engine/modules/registry/pkg"
	uuidpkg "engine/modules/uuid/pkg"
	"engine/services/clock"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type setup struct {
	Assets assets.Service `inject:"1"`

	Registry registry.Service `inject:"1"`
}

func NewSetup() setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		clock.Package(time.RFC3339Nano),
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		registrypkg.Package(),
		ecs.Package(),
		uuidpkg.Package(),
		assetspkg.Package(""),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	return ioc.GetServices[setup](c)
}
