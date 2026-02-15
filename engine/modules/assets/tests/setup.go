package test

import (
	"engine/modules/assets"
	assetspkg "engine/modules/assets/pkg"
	"engine/services/clock"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type setup struct {
	Extensions assets.Extensions `inject:"1"`
	Assets     assets.Service    `inject:"1"`
}

func NewSetup() setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		clock.Package(time.RFC3339Nano),
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		assetspkg.Package(""),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	return ioc.GetServices[setup](c)
}
